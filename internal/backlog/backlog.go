// Package backlog parses BACKLOG.md into structured items so tooling can keep
// GitHub issues and milestones in sync with it. BACKLOG.md is the single source
// of truth (stable TT-n ids); this package never writes it.
package backlog

import (
	"regexp"
	"strings"
)

// Item is one backlog entry (a "TT-n" todo).
type Item struct {
	ID          string // e.g. "TT-4"
	Title       string // issue title, e.g. "TT-4 — Release GitHub sui tag"
	Description string // text after the bold title
	Milestone   string // cleaned section heading, e.g. "M1 — Tooling & release"
	Done        bool   // true when the checkbox is [x]
}

// Milestone is a section of the backlog with its completion state.
type Milestone struct {
	Title string
	Total int
	Done  int
}

// Complete reports whether every item of the milestone is checked. An empty
// milestone is not complete: closing a milestone with no items would hide a
// section that was just created and not yet filled in.
func (m Milestone) Complete() bool { return m.Total > 0 && m.Done == m.Total }

// itemRe matches a checklist line like:
//
//   - [ ] **TT-4 — Release GitHub sui tag**: workflow che ...
var itemRe = regexp.MustCompile(`^\s*-\s*\[([ xX])\]\s*\*\*(TT-\d+)\s*—\s*(.*?)\*\*\s*:?\s*(.*)$`)

var headingRe = regexp.MustCompile(`^##\s+(.*)$`)

// Parse extracts the TT-n items from BACKLOG.md content, in file order.
func Parse(md string) []Item {
	var items []Item
	milestone := ""
	for _, line := range strings.Split(md, "\n") {
		if m := headingRe.FindStringSubmatch(line); m != nil {
			milestone = cleanMilestone(m[1])
			continue
		}
		m := itemRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name := stripBackticks(strings.TrimSpace(m[3]))
		items = append(items, Item{
			ID:          m[2],
			Title:       m[2] + " — " + name,
			Description: strings.TrimSpace(m[4]),
			Milestone:   milestone,
			Done:        m[1] == "x" || m[1] == "X",
		})
	}
	return items
}

// Milestones groups items by milestone, in first-seen order, so callers can
// decide which GitHub milestones to open or close. Items with no milestone
// (before the first heading) are ignored.
func Milestones(items []Item) []Milestone {
	var out []Milestone
	index := map[string]int{}
	for _, it := range items {
		if it.Milestone == "" {
			continue
		}
		i, ok := index[it.Milestone]
		if !ok {
			i = len(out)
			index[it.Milestone] = i
			out = append(out, Milestone{Title: it.Milestone})
		}
		out[i].Total++
		if it.Done {
			out[i].Done++
		}
	}
	return out
}

// cleanMilestone trims a heading to its short title: everything before the
// first " (" (which starts the "(~vX)" / comment tail).
func cleanMilestone(h string) string {
	if i := strings.Index(h, " ("); i >= 0 {
		h = h[:i]
	}
	return strings.TrimSpace(h)
}

func stripBackticks(s string) string { return strings.ReplaceAll(s, "`", "") }
