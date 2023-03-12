package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/luyanakat/booking-app/internal/driver"
	"github.com/luyanakat/booking-app/internal/helpers"
	"github.com/luyanakat/booking-app/internal/repository"
	"github.com/luyanakat/booking-app/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/luyanakat/booking-app/internal/config"
	"github.com/luyanakat/booking-app/internal/forms"
	"github.com/luyanakat/booking-app/internal/models"
	"github.com/luyanakat/booking-app/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingsRepo(a),
	}
}

func NewHandle(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "about.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(res.RoomID)
	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("02-01-2006")
	ed := res.EndDate.Format("02-01-2006")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	data := make(map[string]interface{})
	data["reservation"] = res

	render.WriteTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation Handle the posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Phone = r.Form.Get("phone")
	res.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = res

		render.WriteTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(res)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert reser into db")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		RoomID:        res.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send mail
	htmlMsg := fmt.Sprintf(`<strong>Reservation Confirmation</strong> <br>
       Dear %s, <br>
       This is confirm reservation from %s to %s.`,
		res.FirstName, res.StartDate.Format("02-01-2006"), res.EndDate.Format("02-01-2006"))

	msg := models.MailData{
		To:       res.Email,
		From:     "nai@proton.me",
		Subject:  "Reservation Confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	htmlMsgOwner := fmt.Sprintf(`<strong>You have new reservation</strong> <br>
       Dear %s, <br>
       You have new reservation, room is: %s, from %s to %s.`,
		"Owner", res.Room.RoomName, res.StartDate.Format("02-01-2006"),
		res.EndDate.Format("02-01-2006"))

	msgOwner := models.MailData{
		To:       "owner@gmail.com",
		From:     "nai@gmail.com",
		Subject:  "New Reservation",
		Content:  htmlMsgOwner,
		Template: "basic.html",
	}

	m.App.MailChan <- msgOwner

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Normal(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "normal.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Luxurious(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "luxurious.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "02-01-2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "There are no availability! Please choose another date")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.WriteTemplate(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "02-01-2006"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can not get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("02-01-2006")
	ed := reservation.EndDate.Format("02-01-2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.WriteTemplate(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
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

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	RoomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	startDate := r.URL.Query().Get("s")
	endDate := r.URL.Query().Get("e")

	layout := "02-01-2006"
	sd, err := time.Parse(layout, startDate)
	if err != nil {
		helpers.ServerError(w, err)
	}
	ed, err := time.Parse(layout, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	room, err := m.DB.GetRoomByID(RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation
	res.Room.RoomName = room.RoomName
	res.RoomID = RoomID
	res.StartDate = sd
	res.EndDate = ed

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		render.WriteTemplate(w, r, "login.page.gohtml", &models.TemplateData{})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	user, err := m.DB.GetUserByID(id)
	if err != nil {
		log.Println(err)
	}

	userName := user.LastName

	m.App.Session.Put(r.Context(), "user_name", userName)
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) LogOut(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	m.App.Session.Put(r.Context(), "flash", "Log out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) AdminDashBoard(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "admin-dashboard.page.gohtml", &models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.WriteTemplate(w, r, "admin-new-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var sd, ed string

	for _, reservation := range reservations {
		sd = reservation.StartDate.Format("02-01-2006")
		ed = reservation.EndDate.Format("02-01-2006")
	}

	layout := "02-01-2006"
	for _, reservation := range reservations {
		reservation.StartDate, _ = time.Parse(layout, sd)
		reservation.EndDate, _ = time.Parse(layout, ed)
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.WriteTemplate(w, r, "admin-all-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.WriteTemplate(w, r, "admin-reservation-show.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMapValue := make(map[string]string)
	stringMapValue["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	log.Println(src)

	m.App.Session.Put(r.Context(), "flash", "Update successful!")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminCalendar(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "admin-calendar.page.gohtml", &models.TemplateData{})
}
