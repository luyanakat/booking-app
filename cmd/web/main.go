package main

import (
	"fmt"
	"github.com/luyanakat/booking-app/internal/config"
	"github.com/luyanakat/booking-app/internal/handlers"
	"github.com/luyanakat/booking-app/internal/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// In deploy mode, change this
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true /////
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandle(repo)

	render.NewTemplates(&app)

	fmt.Printf("Start server on port: %s \n", portNumber)
	server := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)

}
