package handlers

import (
	"net/http"

	"github.com/luyanakat/booking-app/pkg/config"
	"github.com/luyanakat/booking-app/pkg/models"
	"github.com/luyanakat/booking-app/pkg/render"
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
	render.WriteTemplate(w, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	userName := m.App.Session.GetString(r.Context(), "remote_user")

	stringMap["remote_ip"] = remoteIp
	stringMap["remote_user"] = userName
	stringMap["name"] = "Nai"
	render.WriteTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}