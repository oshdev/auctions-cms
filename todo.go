package todo

import "github.com/google/uuid"

type Todo struct {
	ID   string
	Name string
}

func NewTodo(name string) Todo {
	return Todo{
		Name: name,
		ID:   uuid.New().String(),
	}
}
