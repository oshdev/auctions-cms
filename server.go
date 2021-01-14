package todo

import (
	"fmt"
	"html/template"
	"net/http"
)

type Repo interface {
	GetTodos() []Todo
	AddTodo(name string)
	DeleteTodo(id string)
}

type Server struct {
	todoTemplate *template.Template
	repo         Repo
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
	if r.Method == http.MethodGet {
		t.todoTemplate.Execute(w, t.repo.GetTodos())
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "couldn't parse the form %v", err)
			return
		}

		if r.URL.Path == "/delete" {
			id := r.FormValue("id")
			t.repo.DeleteTodo(id)
		} else {
			t.repo.AddTodo(r.FormValue("new-item"))
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
