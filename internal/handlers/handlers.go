package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wagnojunior/booking/internal/config"
	"github.com/wagnojunior/booking/internal/driver"
	"github.com/wagnojunior/booking/internal/forms"
	"github.com/wagnojunior/booking/internal/helpers"
	"github.com/wagnojunior/booking/internal/models"
	"github.com/wagnojunior/booking/internal/render"
	"github.com/wagnojunior/booking/internal/repository"
	"github.com/wagnojunior/booking/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Render the template
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// Render the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// MakeReservation is the handler for the make reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	// Creates an empty model reservation and stores it the same format as templatedata.Data
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // Pass an empty form to the make-reservation template
		Data: data,           // Pass the empty reservation model to the make-reservation template
	})
}

// PostMakeReservation handles the posting of a reservation form
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	// Parse the form and check for errors
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// Parse dates from string to time format
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Parse from string to int
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// reservation holds the data from the resservation form, which was entered by the user
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	// Creates a form object
	form := forms.New(r.PostForm)

	// Server-side form validation
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	// Check if the form is NOT valid and send back the form data to the make-reservation template
	if !form.Valid() {
		// data has the same format as Data in templatedata.go
		data := make(map[string]interface{})
		data["reservation"] = reservation

		// Render the template again, but now sending the data from the form to be repopulated
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	// Save to database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Creates a restriction with the newly created reservation (ID)
	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	// Insert restriction to database
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Stores the variable reservation in the session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// Redirects to the reservation.summmary page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// PandaSuite is the handler for the Panda Suite  page
func (m *Repository) PandaSuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "panda-suite.page.tmpl", &models.TemplateData{})
}

// BambooDorm is the handler for the Bamboo Dorm page
func (m *Repository) BambooDorm(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "bamboo-dorm.page.tmpl", &models.TemplateData{})
}

// SearchAvailability is the handler for the Book Now page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability is the handler for the Book Now page
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start") // The argument <start> matches the input name in the form in <search-availability.page.tmpl>
	end := r.Form.Get("end")     // The argument <end> matches the input name in the form in <search-availability.page.tmpl>

	// Parse dates from string to time format
	layout := "2006/01/02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	// Add a Reservation Model to the session
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	// Renders the template reservation-summary and passes the session information to it
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Defines a type to represent a json format
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles requests for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// Creates and populates a variable <resp> of type <jsonResponse>
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	// Formats to json format based on the json tags defined within the ``
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact is the handler for the Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Contact is the handler for the Contact page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// Gets the reserevation from the session. The command .(models.Reservation) is called type assertion.
	// It makes sure that the variable stores in the session is of type models.Reservation.
	// If it is, then ok is set to true; otherwise, it is set to false
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// After the value is retrieved from the session, it is recommended to remove it
	m.App.Session.Remove(r.Context(), "reservation")

	// If ok is true, then store reservation in a format that matches the templatedata.Data
	data := make(map[string]interface{})
	data["reservation"] = reservation

	// Renders the template reservation-summary and passes the session information to it
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}
