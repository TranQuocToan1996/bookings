package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/TranQuocToan1996/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	// Log to terminal
	app.InfoLog.Println("Client error with status of: ", status)
	// Response to user
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Log to terminal
	app.ErrorLog.Println(trace)
	// Response to user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticate(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_id")
}
