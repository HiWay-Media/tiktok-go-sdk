// Command backlog-sync keeps GitHub issues and milestones in sync with
// BACKLOG.md.
//
// BACKLOG.md is the single source of truth: every TT-n item becomes an issue
// (label "backlog", grouped by the milestone of its "## Mn — ..." section).
// Checking an item ([x]) closes its issue; unchecking reopens it. A milestone
// whose items are all checked is closed, and reopened as soon as one is not.
// The sync is idempotent — issues are matched by their "TT-n" title prefix — so
// it is safe to run repeatedly (locally or from
// .github/workflows/backlog-sync.yml).
//
//	go run ./cmd/backlog-sync [-backlog BACKLOG.md] [-dry-run]
//
// Requires the `gh` CLI, authenticated (GH_TOKEN in CI).
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/HiWay-Media/tiktok-go-sdk/internal/backlog"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "print the actions without performing them")
	path := flag.String("backlog", "BACKLOG.md", "path to BACKLOG.md")
	flag.Parse()

	raw, err := os.ReadFile(*path)
	if err != nil {
		fatal(err)
	}
	items := backlog.Parse(string(raw))
	if len(items) == 0 {
		fatal(fmt.Errorf("no TT-n item in %s", *path))
	}

	s := &syncer{dryRun: *dryRun}
	if err := s.run(items); err != nil {
		fatal(err)
	}
	fmt.Printf("backlog-sync: %d item · issue: %d create, %d chiuse, %d riaperte, %d rimilestonate, %d invariate · milestone: %d create, %d chiuse, %d riaperte (dry-run=%v)\n",
		len(items), s.created, s.closed, s.reopened, s.remilestoned, s.unchanged,
		s.msCreated, s.msClosed, s.msReopened, *dryRun)
}

type syncer struct {
	dryRun bool

	created, closed, reopened, remilestoned, unchanged int
	msCreated, msClosed, msReopened                    int
}

type ghIssue struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	State     string `json:"state"` // "OPEN" | "CLOSED"
	Milestone *struct {
		Title string `json:"title"`
	} `json:"milestone"`
}

type ghMilestone struct {
	Number int
	Title  string
	State  string // "open" | "closed"
}

func (s *syncer) run(items []backlog.Item) error {
	s.ensureLabel()
	milestones := backlog.Milestones(items)
	have, err := s.ensureMilestones(milestones)
	if err != nil {
		return err
	}
	existing, err := s.listIssues()
	if err != nil {
		return err
	}
	for _, it := range items {
		if err := s.syncItem(it, existing[it.ID]); err != nil {
			return err
		}
	}
	// Milestone states are reconciled last: closing a milestone before its
	// issues are closed would leave GitHub showing a closed milestone with open
	// work in it.
	return s.syncMilestoneStates(milestones, have)
}

func (s *syncer) syncItem(it backlog.Item, cur *ghIssue) error {
	if cur == nil {
		return s.createIssue(it)
	}
	// An item moved to another section must follow it, otherwise the milestone
	// progress bar on GitHub silently stops matching the backlog.
	if it.Milestone != "" && milestoneTitle(cur) != it.Milestone {
		s.remilestoned++
		if err := s.act("rimilestono", it, func() error {
			return ghVoid("issue", "edit", strconv.Itoa(cur.Number), "--milestone", it.Milestone)
		}); err != nil {
			return err
		}
	}
	switch {
	case it.Done && cur.State == "OPEN":
		s.closed++
		return s.act("chiudo", it, func() error {
			return ghVoid("issue", "close", strconv.Itoa(cur.Number), "--comment", "Completata: item spuntato in BACKLOG.md.")
		})
	case !it.Done && cur.State == "CLOSED":
		s.reopened++
		return s.act("riapro", it, func() error { return ghVoid("issue", "reopen", strconv.Itoa(cur.Number)) })
	default:
		s.unchanged++
		return nil
	}
}

func (s *syncer) createIssue(it backlog.Item) error {
	s.created++
	return s.act("creo", it, func() error {
		args := []string{"issue", "create", "--title", it.Title, "--body", issueBody(it), "--label", "backlog"}
		if it.Milestone != "" {
			args = append(args, "--milestone", it.Milestone)
		}
		out, err := gh(args...)
		if err != nil {
			return err
		}
		// A done item is recorded as a closed issue for full history.
		if it.Done {
			if n := lastPathInt(out); n > 0 {
				return ghVoid("issue", "close", strconv.Itoa(n), "--comment", "Item already completed in BACKLOG.md.")
			}
		}
		return nil
	})
}

// act logs the action and runs it unless in dry-run mode.
func (s *syncer) act(verb string, it backlog.Item, do func() error) error {
	fmt.Printf("  %-12s %s (%s)\n", verb, it.ID, it.Milestone)
	if s.dryRun {
		return nil
	}
	return do()
}

func (s *syncer) ensureLabel() {
	// Best-effort: fails harmlessly if the label already exists.
	if s.dryRun {
		return
	}
	_, _ = gh("label", "create", "backlog", "--color", "5319e7", "--description", "BACKLOG.md item (TT-n), synced by backlog-sync")
}

// ensureMilestones creates the milestones missing on GitHub and returns the
// full set, keyed by title.
func (s *syncer) ensureMilestones(milestones []backlog.Milestone) (map[string]ghMilestone, error) {
	have, err := s.listMilestones()
	if err != nil {
		return nil, err
	}
	for _, m := range milestones {
		if _, ok := have[m.Title]; ok {
			continue
		}
		s.msCreated++
		fmt.Printf("  %-12s %s\n", "milestone", m.Title)
		if s.dryRun {
			continue
		}
		// Tolerate a milestone that appeared behind our back (concurrent run,
		// manual creation): the desired state already holds, so it isn't an error.
		if _, err := gh("api", "-X", "POST", "repos/{owner}/{repo}/milestones", "-f", "title="+m.Title); err != nil && !isAlreadyExists(err) {
			return nil, err
		}
		if refreshed, err := s.listMilestones(); err == nil {
			have = refreshed
		}
	}
	return have, nil
}

// syncMilestoneStates closes a milestone whose items are all checked and
// reopens one that has open work again.
func (s *syncer) syncMilestoneStates(milestones []backlog.Milestone, have map[string]ghMilestone) error {
	for _, m := range milestones {
		cur, ok := have[m.Title]
		if !ok {
			// Only possible in dry-run, where the milestone was not created.
			continue
		}
		want := "open"
		if m.Complete() {
			want = "closed"
		}
		if cur.State == want {
			continue
		}
		verb := "riapro ms"
		s.msReopened++
		if want == "closed" {
			verb = "chiudo ms"
			s.msReopened--
			s.msClosed++
		}
		fmt.Printf("  %-12s %s (%d/%d)\n", verb, m.Title, m.Done, m.Total)
		if s.dryRun {
			continue
		}
		if _, err := gh("api", "-X", "PATCH", fmt.Sprintf("repos/{owner}/{repo}/milestones/%d", cur.Number), "-f", "state="+want); err != nil {
			return err
		}
	}
	return nil
}

// listMilestones returns every milestone of the repo, keyed by title.
//
// --paginate is required: GitHub returns 30 milestones per page by default, so
// without it the newest ones fall off the list and we try to re-create them
// (HTTP 422 already_exists). The TSV projection keeps the output stable across
// the array-per-page shape that --paginate produces.
func (s *syncer) listMilestones() (map[string]ghMilestone, error) {
	out, err := gh("api", "--paginate", "repos/{owner}/{repo}/milestones?state=all&per_page=100",
		"--jq", `.[] | [.number, .state, .title] | @tsv`)
	if err != nil {
		return nil, err
	}
	have := map[string]ghMilestone{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			continue
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		have[parts[2]] = ghMilestone{Number: n, State: parts[1], Title: parts[2]}
	}
	return have, nil
}

func (s *syncer) listIssues() (map[string]*ghIssue, error) {
	out, err := gh("issue", "list", "--state", "all", "--limit", "300", "--json", "number,title,state,milestone")
	if err != nil {
		return nil, err
	}
	var issues []ghIssue
	if err := json.Unmarshal([]byte(out), &issues); err != nil {
		return nil, err
	}
	byID := map[string]*ghIssue{}
	for i := range issues {
		if id := idFromTitle(issues[i].Title); id != "" {
			byID[id] = &issues[i]
		}
	}
	return byID, nil
}

func milestoneTitle(i *ghIssue) string {
	if i == nil || i.Milestone == nil {
		return ""
	}
	return i.Milestone.Title
}

func issueBody(it backlog.Item) string {
	return fmt.Sprintf("%s\n\n---\nTracked in `BACKLOG.md` (**%s**) · milestone _%s_.\n\n"+
		"Managed by `cmd/backlog-sync`: edit the **BACKLOG**, not this issue. "+
		"Check the item (`[x]`) and the issue will be closed on the next sync.", it.Description, it.ID, it.Milestone)
}

// idFromTitle returns the leading "TT-n" token of an issue title, or "".
func idFromTitle(title string) string {
	fields := strings.Fields(title)
	if len(fields) > 0 && strings.HasPrefix(fields[0], "TT-") {
		return fields[0]
	}
	return ""
}

// lastPathInt parses the trailing integer of the last URL printed by gh
// (e.g. https://github.com/owner/repo/issues/42 -> 42).
func lastPathInt(s string) int {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "/"); i >= 0 {
		if n, err := strconv.Atoi(strings.TrimSpace(s[i+1:])); err == nil {
			return n
		}
	}
	return 0
}

// isAlreadyExists reports whether err is the GitHub validation failure raised
// when creating a resource whose unique field is taken (HTTP 422).
func isAlreadyExists(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already_exists")
}

func gh(args ...string) (string, error) {
	cmd := exec.Command("gh", args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		// The API error body lands on stdout ("gh: Validation Failed (HTTP 422)"
		// alone on stderr says nothing), so report both.
		msg := strings.TrimSpace(errb.String())
		if body := strings.TrimSpace(out.String()); body != "" {
			msg = strings.TrimSpace(msg + " " + body)
		}
		return "", fmt.Errorf("gh %s: %v: %s", strings.Join(args, " "), err, msg)
	}
	return out.String(), nil
}

func ghVoid(args ...string) error {
	_, err := gh(args...)
	return err
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "backlog-sync:", err)
	os.Exit(1)
}
