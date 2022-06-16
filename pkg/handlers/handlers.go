package handlers

import (
	"net/http"

	"github.com/wagnojunior/booking/pkg/config"
	"github.com/wagnojunior/booking/pkg/models"
	"github.com/wagnojunior/booking/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Get the remote IP from the http request
	remoteIP := r.RemoteAddr

	// Adds <remoteIP> to the session
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	// Render the template
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// Adds a string to the <stringMap>
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	// Gets the remote IP from the session and adds to <stringMap>
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// Render the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// MakeReservation is the handler for the make reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// PandaSuite is the handler for the Panda Suite  page
func (m *Repository) PandaSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "panda-suite.page.tmpl", &models.TemplateData{})
}

// BambooDorm is the handler for the Bamboo Dorm page
func (m *Repository) BambooDorm(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "bamboo-dorm.page.tmpl", &models.TemplateData{})
}

// SearchAvailability is the handler for the Book Now page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Contact is the handler for the Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}
