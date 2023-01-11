package main

import (
	"github.com/luyanakat/booking-app/internal/config"
	"github.com/luyanakat/booking-app/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/normal-room", handlers.Repo.Normal)
	mux.Get("/luxurious-room", handlers.Repo.Luxurious)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/make-reservation", handlers.Repo.Reservation)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
