package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// fileData is the on-disk JSON format.
type fileData struct {
	NextID int    `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}

// Store manages tasks in a JSON file.
type Store struct {
	path string
}

// NewStore creates a store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// DefaultPath returns ~/.tasks.json.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, ".tasks.json"), nil
}

// load reads all data from the file. Returns zero value if file doesn't exist.
func (s *Store) load() (fileData, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fileData{NextID: 1}, nil
		}
		return fileData{}, fmt.Errorf("read %s: %w", s.path, err)
	}
	if len(data) == 0 {
		return fileData{NextID: 1}, nil
	}
	var fd fileData
	if err := json.Unmarshal(data, &fd); err != nil {
		return fileData{}, fmt.Errorf("unmarshal tasks: %w", err)
	}
	if fd.NextID == 0 {
		fd.NextID = 1
	}
	return fd, nil
}

// save writes all data to the file.
func (s *Store) save(fd fileData) error {
	data, err := json.MarshalIndent(fd, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", s.path, err)
	}
	return nil
}

// Add creates a new task with the given title.
func (s *Store) Add(title string) (Task, error) {
	fd, err := s.load()
	if err != nil {
		return Task{}, err
	}
	t := Task{
		ID:        fd.NextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	fd.NextID++
	fd.Tasks = append(fd.Tasks, t)
	if err := s.save(fd); err != nil {
		return Task{}, err
	}
	return t, nil
}

// List returns all tasks, optionally filtered.
// filter: nil = all, ptr to true = done only, ptr to false = pending only.
func (s *Store) List(filter *bool) ([]Task, error) {
	fd, err := s.load()
	if err != nil {
		return nil, err
	}
	if filter == nil {
		return fd.Tasks, nil
	}
	var out []Task
	for _, t := range fd.Tasks {
		if t.Done == *filter {
			out = append(out, t)
		}
	}
	return out, nil
}

// Done marks a task as done by ID.
func (s *Store) Done(id int) error {
	fd, err := s.load()
	if err != nil {
		return err
	}
	for i := range fd.Tasks {
		if fd.Tasks[i].ID == id {
			fd.Tasks[i].Done = true
			return s.save(fd)
		}
	}
	return fmt.Errorf("task %d not found", id)
}

// Remove deletes a task by ID.
func (s *Store) Remove(id int) error {
	fd, err := s.load()
	if err != nil {
		return err
	}
	for i, t := range fd.Tasks {
		if t.ID == id {
			fd.Tasks = append(fd.Tasks[:i], fd.Tasks[i+1:]...)
			return s.save(fd)
		}
	}
	return fmt.Errorf("task %d not found", id)
}
