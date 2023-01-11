package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/luyanakat/booking-app/internal/config"
	"github.com/luyanakat/booking-app/internal/models"
	"github.com/luyanakat/booking-app/internal/render"
	"log"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandle(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	userName := r.UserAgent()

	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	m.App.Session.Put(r.Context(), "remote_user", userName)
	render.WriteTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIp

	render.WriteTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.WriteTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{})
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

	w.Write([]byte(fmt.Sprintf("Start date is %s and the end date is: %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
