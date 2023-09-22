package server

import (
	"log"
	"net/http"
	"path/filepath"

	"html/template"
)

// Provides an HTML page for the admin interface.
// Source: https://www.alexedwards.net/blog/serving-static-sites-with-go

type AdminHandler struct {
	http.Handler
	userStore UserStore

	tmpl *template.Template
}

func NewAdminHandler(store UserStore) *AdminHandler {
	ah := new(AdminHandler)
	ah.userStore = store

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(ah.HandleHtml))
	router.Handle("/add-quiz", http.HandlerFunc(ah.AddQuiz))
	router.Handle("/reset-pwd", http.HandlerFunc(ah.HandleHtml))

	// load HTML template
	layoutPath := filepath.Join("templates", "layout.html")
	contentPath := filepath.Join("templates", "admin.html")
	tmpl, err := template.ParseFiles(layoutPath, contentPath)
	if err != nil {
		panic("Error loading template files")
	}
	ah.tmpl = tmpl

	ah.Handler = router
	return ah
}

func (ah *AdminHandler) HandleHtml(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	statItems, err := ah.userStore.GetStats()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	htmlData := &TemplateData{StatItems: statItems}

	err = ah.tmpl.ExecuteTemplate(w, "admin_page", htmlData)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (ah *AdminHandler) AddQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// generate random quiz
	randomQuizProperties := RandomQuizProperties()
	_, err := ah.userStore.GenerateQuiz(randomQuizProperties, true)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// redirect to admin page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
