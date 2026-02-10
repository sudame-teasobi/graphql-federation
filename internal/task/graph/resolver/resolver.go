package resolver

import "gft/internal/task"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Repo *task.Repository
}
