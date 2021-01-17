package todo

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Repo interface {
	GetTodos() []Todo
	AddTodo(name string)
	DeleteTodo(id string)
	GetTodo(id string) Todo
	EditTodo(id string, name string)
}

func NewServer(templateFolderPath string, repo Repo) (*mux.Router, error) {
	todoTemplate, err := template.ParseGlob(templateFolderPath)
	if err != nil {
		return nil, fmt.Errorf(
			"could not load todo template from %q, %v",
			templateFolderPath,
			err,
		)
	}
	router := mux.NewRouter()

	server := server{
		todoTemplate: todoTemplate,
		repo:         repo,
	}

	router.HandleFunc("/", server.viewTodos).Methods(http.MethodGet)
	router.HandleFunc("/add", server.viewAddTodoForm).Methods(http.MethodGet)
	router.HandleFunc("/add", server.addTodo).Methods(http.MethodPost)
	router.HandleFunc("/delete", server.deleteTodoStreamed).Methods(http.MethodPost).Headers("Accept", "text/vnd.turbo-stream.html")
	router.HandleFunc("/delete", server.deleteTodo).Methods(http.MethodPost)
	router.HandleFunc("/edit/{id}", server.viewEditTodoForm).Methods(http.MethodGet)
	router.HandleFunc("/edit/{id}", server.editTodo).Methods(http.MethodPost)
	return router, nil
}

type server struct {
	todoTemplate *template.Template
	repo         Repo
}

func (s *server) viewTodos(writer http.ResponseWriter, r *http.Request) {
	s.todoTemplate.ExecuteTemplate(writer, "todo.gohtml", s.repo.GetTodos())
}

func (s *server) viewAddTodoForm(writer http.ResponseWriter, r *http.Request) {
	s.todoTemplate.ExecuteTemplate(writer, "add.gohtml", s.repo.GetTodos())
}

func (s *server) addTodo(writer http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	s.repo.AddTodo(r.PostForm.Get("new-item"))
	http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func (s *server) deleteTodoStreamed(writer http.ResponseWriter, r *http.Request) {
	id := getIdFromForm(r)

	todo := s.repo.GetTodo(id)
	s.repo.DeleteTodo(id)
	writer.Header().Add("Content-Type", "text/vnd.turbo-stream.html")

	s.todoTemplate.ExecuteTemplate(writer, "toaster.partial.gohtml", todo)
	s.todoTemplate.ExecuteTemplate(writer, "replace-todo-list-stream.gohtml", s.repo.GetTodos())
}

func (s *server) deleteTodo(writer http.ResponseWriter, r *http.Request) {
	s.repo.DeleteTodo(getIdFromForm(r))
	http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func (s *server) viewEditTodoForm(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	s.todoTemplate.ExecuteTemplate(writer, "edit.gohtml", s.repo.GetTodo(id))
}

func (s *server) editTodo(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseMultipartForm(1024)

	id := vars["id"]
	s.repo.EditTodo(id, r.PostForm.Get("updated-name"))
	http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func getIdFromForm(r *http.Request) string {
	r.ParseMultipartForm(1024)
	id := r.PostForm.Get("id")
	return id
}
