package handlers

import (
	"encoding/json"
	"errors"
	"log"
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	// get the room information by ID
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Add the room name to the reservation model
	res.Room.RoomName = room.RoomName

	// Add the reservation model <res> to the session
	m.App.Session.Put(r.Context(), "reservation", res)

	// Format the date to string
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil), // Pass an empty form to the make-reservation template
		Data:      data,           // Pass the empty reservation model to the make-reservation template
		StringMap: stringMap,
	})
}

// PostMakeReservation handles the posting of a reservation form
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get from session"))
		return
	}

	// Parse the form and check for errors
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

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
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
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
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles requests for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	// Get the fields BY NAME from the form
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	// Convert sd and ed from string to date
	layout := "2006/01/02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	log.Println(startDate, endDate)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	// Creates and populates a variable <resp> of type <jsonResponse>
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	// Renders the template reservation-summary and passes the session information to it
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
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

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make a reservation
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// We have to grab the values from the URL (id, s, e)
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id")) // Get the ID from the Request and convert it to str
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	// Convert sd and ed from string to date
	layout := "2006/01/02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	// Create a variable of type <Reservation>
	var res models.Reservation

	// get the room name by ID
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Populate the variable <res>
	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	// Put <res> into the session
	m.App.Session.Put(r.Context(), "reservation", res)

	// Redirect
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
