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

type Server struct {
	todoTemplate *template.Template
	repo         Repo
	router       *mux.Router
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
	router := mux.NewRouter()

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		todoTemplate.ExecuteTemplate(writer, "todo.gohtml", repo.GetTodos())
	}).Methods(http.MethodGet)

	router.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		todoTemplate.ExecuteTemplate(writer, "add.gohtml", repo.GetTodos())
	}).Methods(http.MethodGet)

	router.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(1024)
		repo.AddTodo(request.PostForm.Get("new-item"))
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}).Methods(http.MethodPost)

	router.HandleFunc("/delete", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(1024)
		repo.DeleteTodo(request.PostForm.Get("id"))
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}).Methods(http.MethodPost)

	router.HandleFunc("/edit/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]
		todoTemplate.ExecuteTemplate(writer, "edit.gohtml", repo.GetTodo(id))
	}).Methods(http.MethodGet)

	router.HandleFunc("/edit/{id}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		request.ParseMultipartForm(1024)

		id := vars["id"]
		repo.EditTodo(id, request.PostForm.Get("updated-name"))
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}).Methods(http.MethodPost)

	return &Server{
		todoTemplate: todoTemplate,
		repo:         repo,
		router:       router,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
