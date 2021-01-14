package todo

import "github.com/google/uuid"

type Todo struct {
	ID string
	Name string
	Completed bool
}

func NewTodo(name string) Todo {
	return Todo{
		Name: name,
		ID: uuid.New().String(),
		Completed: false,
	}
}
