package todo

import (
	"fmt"
	"html/template"
	"net/http"
)

type Repo interface {
	GetTodos() []Todo
	AddTodo(name string)
	ToggleTodoComplete(id string)
	DeleteTodo(id string)
}

type Server struct {
	todoTemplate *template.Template
	repo Repo
}

func NewServer(templateFolderPath string, repo Repo) (*Server, error) {
	todoTemplate, err := template.ParseGlob(templateFolderPath)
	if err != nil {
		return nil, fmt.Errorf(
			"could not load todo template from %q, %v",
			templateFolderPath,
			err,
		)
	}
	return &Server{todoTemplate: todoTemplate, repo: repo}, nil
}

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method==http.MethodGet {
		t.todoTemplate.Execute(w, t.repo.GetTodos())
	}
}
