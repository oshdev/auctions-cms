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
	EditTodo(id string, newName string)
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

	server := server{
		todoTemplate: todoTemplate,
		repo:         repo,
	}

	router := mux.NewRouter()
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

func (s *server) viewTodos(w http.ResponseWriter, r *http.Request) {
	s.todoTemplate.ExecuteTemplate(w, "todo.gohtml", s.repo.GetTodos())
}

func (s *server) viewAddTodoForm(w http.ResponseWriter, r *http.Request) {
	s.todoTemplate.ExecuteTemplate(w, "add.gohtml", s.repo.GetTodos())
}

func (s *server) addTodo(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	s.repo.AddTodo(r.PostForm.Get("new-item"))
	redirectToHome(w, r)
}

func (s *server) deleteTodoStreamed(w http.ResponseWriter, r *http.Request) {
	id := getIdFromForm(r)

	todo := s.repo.GetTodo(id)
	s.repo.DeleteTodo(id)
	w.Header().Add("Content-Type", "text/vnd.turbo-stream.html")

	s.todoTemplate.ExecuteTemplate(w, "toaster.partial.gohtml", todo)
	s.todoTemplate.ExecuteTemplate(w, "replace-todo-list-stream.gohtml", s.repo.GetTodos())
}

func (s *server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	s.repo.DeleteTodo(getIdFromForm(r))
	redirectToHome(w, r)
}

func (s *server) viewEditTodoForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	s.todoTemplate.ExecuteTemplate(w, "edit.gohtml", s.repo.GetTodo(id))
}

func (s *server) editTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseMultipartForm(1024)

	id := vars["id"]
	s.repo.EditTodo(id, r.PostForm.Get("updated-name"))
	redirectToHome(w, r)
}

func redirectToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getIdFromForm(r *http.Request) string {
	r.ParseMultipartForm(1024)
	id := r.PostForm.Get("id")
	return id
}
