package task

import (
	"os"
	"path/filepath"
	"testing"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	p := filepath.Join(t.TempDir(), "tasks.json")
	return NewStore(p)
}

func TestAdd(t *testing.T) {
	s := tempStore(t)

	tests := []struct {
		title  string
		wantID int
	}{
		{"buy milk", 1},
		{"walk dog", 2},
		{"read book", 3},
	}
	for _, tt := range tests {
		task, err := s.Add(tt.title)
		if err != nil {
			t.Fatalf("add %q: %v", tt.title, err)
		}
		if task.ID != tt.wantID {
			t.Errorf("add %q: got ID %d, want %d", tt.title, task.ID, tt.wantID)
		}
		if task.Title != tt.title {
			t.Errorf("got title %q, want %q", task.Title, tt.title)
		}
		if task.Done {
			t.Errorf("new task should not be done")
		}
		if task.CreatedAt.IsZero() {
			t.Errorf("CreatedAt should be set")
		}
	}
}

func TestList(t *testing.T) {
	s := tempStore(t)

	// empty store
	tasks, err := s.List(nil)
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}

	s.Add("one")
	s.Add("two")
	s.Done(1)

	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name   string
		filter *bool
		want   int
	}{
		{"all", nil, 2},
		{"done", boolPtr(true), 1},
		{"pending", boolPtr(false), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.List(tt.filter)
			if err != nil {
				t.Fatalf("list: %v", err)
			}
			if len(got) != tt.want {
				t.Errorf("got %d tasks, want %d", len(got), tt.want)
			}
		})
	}
}

func TestDone(t *testing.T) {
	s := tempStore(t)
	s.Add("task one")

	if err := s.Done(1); err != nil {
		t.Fatalf("done: %v", err)
	}

	tasks, _ := s.List(nil)
	if !tasks[0].Done {
		t.Error("task should be marked done")
	}

	// not found
	if err := s.Done(99); err == nil {
		t.Error("expected error for missing task")
	}
}

func TestRemove(t *testing.T) {
	s := tempStore(t)
	s.Add("task one")
	s.Add("task two")

	if err := s.Remove(1); err != nil {
		t.Fatalf("remove: %v", err)
	}

	tasks, _ := s.List(nil)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != 2 {
		t.Errorf("remaining task ID should be 2, got %d", tasks[0].ID)
	}

	// not found
	if err := s.Remove(99); err == nil {
		t.Error("expected error for missing task")
	}
}

func TestRoundTrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "tasks.json")

	// write with one store instance
	s1 := NewStore(p)
	s1.Add("persisted task")
	s1.Add("another task")
	s1.Done(2)

	// read with a fresh store instance
	s2 := NewStore(p)
	tasks, err := s2.List(nil)
	if err != nil {
		t.Fatalf("load from new store: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "persisted task" {
		t.Errorf("got title %q, want %q", tasks[0].Title, "persisted task")
	}
	if tasks[1].Done != true {
		t.Error("second task should be done")
	}
}

func TestAutoIncrementAfterRemove(t *testing.T) {
	s := tempStore(t)
	s.Add("one")   // ID 1
	s.Add("two")   // ID 2
	s.Remove(2)    // remove highest
	task, _ := s.Add("three") // should be ID 3, not 2
	if task.ID != 3 {
		t.Errorf("got ID %d after remove, want 3", task.ID)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	s := NewStore(filepath.Join(t.TempDir(), "nope.json"))
	tasks, err := s.List(nil)
	if err != nil {
		t.Fatalf("list on missing file: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestLoadEmptyFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "empty.json")
	os.WriteFile(p, []byte{}, 0644)

	s := NewStore(p)
	tasks, err := s.List(nil)
	if err != nil {
		t.Fatalf("list on empty file: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}
