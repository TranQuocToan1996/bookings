package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteToConsole log some text to terminal when client load a page
func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page: " + r.Host + r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// NoSurf adds CSRF protestion for all POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	// Using cookies to make sure the csrfToken available on a per page basic
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,             // Decrease chance  risk of client side script accessing the protected cookie (stop Js in client side)
		Path:     "/",              // Entire site
		Secure:   app.InProduction, // HTTPs
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// LoadAndSave provides middleware which automatically loads and saves session
// data for the current request, and communicates the session token to and from
// the client in a cookie.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
