package user

import (
	"cmp"
	"fmt"
	"gft/internal/user/model"
	"maps"
	"slices"
)

type Repository struct {
	users map[string]*model.User
}

func NewRepository(initUsers []*model.User) *Repository {
	users := map[string]*model.User{}
	for _, u := range initUsers {
		users[u.ID] = u
	}

	return &Repository{users: users}
}

func (r *Repository) GetByID(id string) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("failed to find user id: %s", id)
	}

	return u, nil
}

func (r *Repository) GetAll() []*model.User {
	return slices.SortedFunc(maps.Values(r.users), func(a, b *model.User) int { return cmp.Compare(a.ID, b.ID) })
}
