package model

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	TaskIds []string `json:"taskIds"`
}

func (User) IsNode()         {}
func (u User) GetID() string { return u.ID }

func (User) IsEntity() {}
