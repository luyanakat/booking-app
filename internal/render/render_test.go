package render

import (
	"github.com/luyanakat/booking-app/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestWriteTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = WriteTemplate(&ww, r, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("error writing templates to browser")
	}

	err = WriteTemplate(&ww, r, "none.page.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("rendered templates does not exist")
	}

	app.UseCache = true
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx, _ := session.Load(r.Context(), r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
