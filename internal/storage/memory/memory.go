package memory

import (
	"sort"
	"sync"
	"time"

	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Store struct {
	mu     sync.RWMutex
	nextID int
	tasks  map[int]models.Task
}

var _ storage.Storage = (*Store)(nil)

func New() *Store { return &Store{nextID: 1, tasks: make(map[int]models.Task)} }

func (s *Store) List() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return tasks
}

func (s *Store) Create(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task.ID = s.nextID
	s.nextID++
	task.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	s.tasks[task.ID] = task
	return task, nil
}

func (s *Store) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *Store) Update(id int, task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.tasks[id]
	if !ok {
		return models.Task{}, storage.ErrNotFound
	}
	task.ID, task.CreatedAt = id, old.CreatedAt
	s.tasks[id] = task
	return task, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return storage.ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
