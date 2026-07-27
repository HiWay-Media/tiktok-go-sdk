package backlog

import (
	"os"
	"strings"
	"testing"
)

const sample = `# Backlog

Intro paragraph, ignored.

## M1 — Tooling & release (~v0.3)

- [x] **TT-1 — Regole di lavoro**: ` + "`CLAUDE.md`" + ` + AGENTS.md. _(v0.3.0)_
- [ ] **TT-5 — README corretto**: l'import punta alla radice.

## M2 — OAuth utente (~v0.4)

- [ ] **TT-7 — Scambio authorization code**: grant authorization_code.
- not a checklist line
- [ ] plain todo without an id
`

func TestParse(t *testing.T) {
	items := Parse(sample)
	if len(items) != 3 {
		t.Fatalf("got %d items, want 3: %+v", len(items), items)
	}

	first := items[0]
	if first.ID != "TT-1" {
		t.Errorf("ID = %q, want TT-1", first.ID)
	}
	if !first.Done {
		t.Error("TT-1 should be done")
	}
	// Backticks are stripped from the title: they would render literally in the
	// GitHub issue title, where markdown is not interpreted.
	if first.Title != "TT-1 — Regole di lavoro" {
		t.Errorf("Title = %q", first.Title)
	}
	if strings.Contains(first.Description, "CLAUDE.md") == false {
		t.Errorf("Description = %q, want the text after the bold title", first.Description)
	}
	if first.Milestone != "M1 — Tooling & release" {
		t.Errorf("Milestone = %q, want the heading without the (~v) tail", first.Milestone)
	}

	if items[1].ID != "TT-5" || items[1].Done {
		t.Errorf("second item = %+v, want open TT-5", items[1])
	}
	if items[2].Milestone != "M2 — OAuth utente" {
		t.Errorf("third item milestone = %q", items[2].Milestone)
	}
}

func TestMilestones(t *testing.T) {
	ms := Milestones(Parse(sample))
	if len(ms) != 2 {
		t.Fatalf("got %d milestones, want 2: %+v", len(ms), ms)
	}
	// First-seen order, so the sync reports milestones in roadmap order.
	if ms[0].Title != "M1 — Tooling & release" || ms[1].Title != "M2 — OAuth utente" {
		t.Fatalf("milestones out of order: %+v", ms)
	}
	if ms[0].Total != 2 || ms[0].Done != 1 || ms[0].Complete() {
		t.Errorf("M1 = %+v, want 1/2 and not complete", ms[0])
	}
	if ms[1].Total != 1 || ms[1].Done != 0 {
		t.Errorf("M2 = %+v, want 0/1", ms[1])
	}
}

func TestMilestoneComplete(t *testing.T) {
	if (Milestone{Title: "M9", Total: 3, Done: 3}).Complete() != true {
		t.Error("a fully checked milestone must be complete")
	}
	// An empty milestone must NOT be complete: a section just added to the
	// roadmap would otherwise be created and immediately closed on GitHub.
	if (Milestone{Title: "M9"}).Complete() {
		t.Error("an empty milestone must not be complete")
	}
}

// TestRealBacklog guards the file the automation actually reads: a typo in an
// id or a heading would silently drop items from the sync.
func TestRealBacklog(t *testing.T) {
	raw, err := os.ReadFile("../../BACKLOG.md")
	if err != nil {
		t.Fatal(err)
	}
	items := Parse(string(raw))
	if len(items) == 0 {
		t.Fatal("no TT-n item parsed from BACKLOG.md")
	}
	seen := map[string]bool{}
	for _, it := range items {
		if seen[it.ID] {
			t.Errorf("duplicate id %s: ids must be unique, they key the issues", it.ID)
		}
		seen[it.ID] = true
		if it.Milestone == "" {
			t.Errorf("%s has no milestone: every item must sit under a '## Mn — ...' heading", it.ID)
		}
		if it.Description == "" {
			t.Errorf("%s has an empty description: it becomes the issue body", it.ID)
		}
	}
	// Every checklist line in the file must be a well-formed item; a line that
	// looks like a todo but does not parse would never reach GitHub.
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- [") {
			continue
		}
		if itemRe.FindStringSubmatch(line) == nil {
			t.Errorf("checklist line does not match the item format: %q", trimmed)
		}
	}
}
