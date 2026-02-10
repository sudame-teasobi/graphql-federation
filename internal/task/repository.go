package task

import (
	"cmp"
	"fmt"
	"gft/internal/task/model"
	"maps"
	"slices"
)

type Repository struct {
	tasks map[string]*model.Task
}

func NewRepository(initTasks []*model.Task) *Repository {
	tasks := map[string]*model.Task{}
	for _, t := range initTasks {
		tasks[t.ID] = t
	}
	return &Repository{tasks: tasks}
}

func (r *Repository) GetByID(id string) (*model.Task, error) {
	t, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("failed to find a task by id: %s", id)
	}
	return t, nil
}

func (r *Repository) GetAll() []*model.Task {
	return slices.SortedFunc(maps.Values(r.tasks), func(a, b *model.Task) int { return cmp.Compare(a.ID, b.ID) })
}
