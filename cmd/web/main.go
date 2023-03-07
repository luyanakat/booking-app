package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/luyanakat/booking-app/internal/config"
	"github.com/luyanakat/booking-app/internal/driver"
	"github.com/luyanakat/booking-app/internal/handlers"
	"github.com/luyanakat/booking-app/internal/helpers"
	"github.com/luyanakat/booking-app/internal/models"
	"github.com/luyanakat/booking-app/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":3000"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Printf("Start server on port: %s \n", portNumber)
	server := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//put in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	// In deploy mode, change this
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true /////
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to db
	log.Println("Connecting to db...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=abc@123")
	if err != nil {
		log.Fatal("Can't connect to database")
	}
	log.Println("Connected to db!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandle(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
