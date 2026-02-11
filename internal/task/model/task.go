package model

type Task struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	UserID string `json:"userId"`
}

func (Task) IsNode()         {}
func (t Task) GetID() string { return t.ID }

func (Task) IsEntity() {}
