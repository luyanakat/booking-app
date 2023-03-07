package helpers

import (
	"fmt"
	"github.com/luyanakat/booking-app/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, statusCode int) {
	app.InfoLog.Println("Client error with status of", statusCode)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")

	return exists

}
