package auction

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Repo interface {
	GetAuctions() []Auction
	AddAuction(name string, seller string, bidder string)
	DeleteAuction(id string)
	GetAuction(id string) Auction
	EditAuction(id string, newName string, seller string, bidder string)
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
	router.HandleFunc("/", server.viewAuctions).Methods(http.MethodGet)
	router.HandleFunc("/add", server.viewAddAuctionForm).Methods(http.MethodGet)
	router.HandleFunc("/add", server.addAuction).Methods(http.MethodPost)
	router.HandleFunc("/delete", server.deleteAuctionStreamed).Methods(http.MethodPost).Headers("Accept", "text/vnd.turbo-stream.html")
	router.HandleFunc("/delete", server.deleteAuction).Methods(http.MethodPost)
	router.HandleFunc("/edit/{id}", server.viewEditAuctionForm).Methods(http.MethodGet)
	router.HandleFunc("/edit/{id}", server.editAuction).Methods(http.MethodPost)
	return router, nil
}

type server struct {
	todoTemplate *template.Template
	repo         Repo
}

func (s *server) viewAuctions(w http.ResponseWriter, r *http.Request) {
	s.todoTemplate.ExecuteTemplate(w, "index.gohtml", s.repo.GetAuctions())
}

func (s *server) viewAddAuctionForm(w http.ResponseWriter, r *http.Request) {
	s.todoTemplate.ExecuteTemplate(w, "add.gohtml", s.repo.GetAuctions())
}

func (s *server) addAuction(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	s.repo.AddAuction(r.PostForm.Get("new-item"), r.PostForm.Get("seller"), r.PostForm.Get("bidder"))
	redirectToHome(w, r)
}

func (s *server) deleteAuctionStreamed(w http.ResponseWriter, r *http.Request) {
	id := getIdFromForm(r)

	todo := s.repo.GetAuction(id)
	s.repo.DeleteAuction(id)
	w.Header().Add("Content-Type", "text/vnd.turbo-stream.html")

	s.todoTemplate.ExecuteTemplate(w, "toaster.partial.gohtml", todo)
	s.todoTemplate.ExecuteTemplate(w, "replace-auction-list-stream.gohtml", s.repo.GetAuctions())
}

func (s *server) deleteAuction(w http.ResponseWriter, r *http.Request) {
	s.repo.DeleteAuction(getIdFromForm(r))
	redirectToHome(w, r)
}

func (s *server) viewEditAuctionForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	s.todoTemplate.ExecuteTemplate(w, "edit.gohtml", s.repo.GetAuction(id))
}

func (s *server) editAuction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseMultipartForm(1024)

	id := vars["id"]
	s.repo.EditAuction(id, r.PostForm.Get("updated-name"), r.PostForm.Get("updated-seller"), r.PostForm.Get("updated-bidder"))
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
