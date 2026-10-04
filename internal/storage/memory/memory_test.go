package memory

import (
	"sync"
	"tasks-api/internal/models"
	"testing"
	"time"
)

func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task, err := s.Create(models.Task{Title: "x"})
			if err != nil {
				t.Error(err)
				return
			}
			if _, ok := s.Get(task.ID); !ok {
				t.Error("created task missing")
			}
			if _, err := s.Update(task.ID, models.Task{Title: "updated", Done: true}); err != nil {
				t.Error(err)
			}
			s.List()
		}()
	}
	wg.Wait()
	list := s.List()
	if len(list) != 100 {
		t.Fatalf("got %d tasks", len(list))
	}
	for i, task := range list {
		if task.ID != i+1 || !task.Done {
			t.Fatalf("unexpected task: %+v", task)
		}
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if err := s.Delete(id); err != nil {
				t.Error(err)
			}
			s.List()
		}(task.ID)
	}
	wg.Wait()
	if len(s.List()) != 0 {
		t.Fatal("delete failed")
	}
}

func TestMetadataAndCopies(t *testing.T) {
	s := New()
	a, _ := s.Create(models.Task{ID: 999, Title: "original", CreatedAt: "invalid"})
	if a.ID != 1 {
		t.Fatal("ID must be generated")
	}
	if _, err := time.Parse(time.RFC3339, a.CreatedAt); err != nil {
		t.Fatal(err)
	}
	list := s.List()
	list[0].Title = "mutated"
	if got, _ := s.Get(a.ID); got.Title != "original" {
		t.Fatal("list leaks mutable state")
	}
	b, _ := s.Update(a.ID, models.Task{ID: 9, Title: "new", CreatedAt: "invalid"})
	if b.ID != a.ID || b.CreatedAt != a.CreatedAt {
		t.Fatal("update changed metadata")
	}
	_ = s.Delete(a.ID)
	c, _ := s.Create(models.Task{Title: "next"})
	if c.ID <= a.ID {
		t.Fatal("ID reused")
	}
}
