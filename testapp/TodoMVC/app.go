package main

import wd "github.com/phaikawl/wade"

// the different states a TodoEntry can be in
const (
	stateEditing   = "editing"
	stateCompleted = "completed"
)

// TodoEntry represents a single entry in the todo list
type TodoEntry struct {
	Text  string
	Done  bool
	State string
	pView *TodoView // pointer to parent view
}

type todoEntryTag struct {
	Entry *TodoEntry
}

// ToggleEdit updates the state for the TodoEntry
func (t *TodoEntry) ToggleEdit() {
	if t.State == stateEditing {
		t.setCompleteState()
	} else {
		t.State = stateEditing
	}
}

// Destroy removes the entry from the list
func (t *TodoEntry) Destroy() {
	println("clicked Destroy:" + t.Text)

	var i int
	var entry *TodoEntry
	for i, entry = range t.pView.Entries {
		if entry.Text == t.Text {
			println("Deleting: " + entry.Text)
			break
		}
	}

	t.pView.DeleteByIndex(i)

}

// ToggleDone switches the Done field on or off
func (t *TodoEntry) ToggleDone() {
	println("clicked ToggleDone:" + t.Text)
	t.Done = !t.Done
	t.setCompleteState()
}

// setCompleteState is just a small helper to reuse this if
func (t *TodoEntry) setCompleteState() {
	if t.Done {
		t.State = stateCompleted
	} else {
		t.State = ""
	}
}

type TodoView struct {
	NewEntry string
	Entries  []*TodoEntry
}

//
func (t *TodoView) ToggleAll() {
	println("clicked ToggleAll")
	for _, e := range t.Entries {
		e.ToggleDone()
	}
}

func (t *TodoView) AddEntry() {
	if t.NewEntry != "" {
		println("Adding:'" + t.NewEntry + "'")
		t.Entries = append(t.Entries, &TodoEntry{Text: t.NewEntry})
		t.NewEntry = ""
	}
}

func (t *TodoView) DeleteByIndex(i int) {
	copy(t.Entries[i:], t.Entries[i+1:])
	t.Entries[len(t.Entries)-1] = nil
	t.Entries = t.Entries[:len(t.Entries)-1]
}

func main() {
	wadeApp := wd.WadeUp("pg-main", "/todo", func(wade *wd.Wade) {
		wade.Pager().RegisterPages("wpage-root")

		// our custom tags
		wade.Custags().RegisterNew("todoentry", "t-todoentry", todoEntryTag{})

		// our main controller
		wade.Pager().RegisterController("pg-main", func(p *wd.PageData) interface{} {
			println("called RegisterController for pg-main")
			view := new(TodoView)

			view.Entries = []*TodoEntry{
				&TodoEntry{pView: view, Text: "create a datastore for entries", Done: true},
				&TodoEntry{pView: view, Text: "add new entries"},
				&TodoEntry{pView: view, Text: "toggle edit off - click anywhere else"},
				&TodoEntry{pView: view, Text: "ToggleAll should do something", Done: true},
				&TodoEntry{pView: view, Text: "destroy -> delete from the list"},
				&TodoEntry{pView: view, Text: "add filters for state"},
				&TodoEntry{pView: view, Text: "update counters in footer"},
			}

			// update the t.State
			// might be better to bind to Done directly
			for _, e := range view.Entries {
				e.setCompleteState()
			}
			return view
		})
	})

	wadeApp.Start()
}
