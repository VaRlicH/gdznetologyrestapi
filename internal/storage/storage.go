package storage

import (
	"errors"
	"tasks-api/internal/models"
)

var ErrNotFound = errors.New("task not found")

type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}
