package todo

type InMemoryRepo struct {
	todos []Todo
}

func NewInMemoryRepo() *InMemoryRepo {
	var todos []Todo
	todos = append(todos, NewTodo("Implement Hotwire"))
	return &InMemoryRepo{todos: todos}
}

func (i *InMemoryRepo) GetTodos() []Todo {
	return i.todos
}

func (i *InMemoryRepo) AddTodo(name string) {
	i.todos = append(i.todos, NewTodo(name))
}

func (i *InMemoryRepo) DeleteTodo(id string) {
	var newList []Todo
	for _, todo := range i.todos {
		if todo.ID != id {
			newList = append(newList, todo)
		}
	}
	i.todos = newList
}
