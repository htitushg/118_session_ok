package main

import (
	"118_session_ok/assets"
	"118_session_ok/controllers"
	"fmt"
	"log"
	"net/http"
)

type apiHandler struct{}

//func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	//mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		/* newSessionToken := uuid.NewString() */
		sessionToken := controllers.Pseudo_uuid()
		maxAge := 120

		// Set the token in the session map, along with the user whom it represents
		assets.Sessions[sessionToken] = assets.Session{
			Pseudo: assets.Sessions[sessionToken].Pseudo,
			MaxAge: maxAge,
		}
		// Set the new token as the users `session_token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  sessionToken,
			MaxAge: maxAge,
		})
		c, err := req.Cookie("session_token")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "%v,\nreq.URL.Path= %v\n", c, req.URL.Path)
		fmt.Fprintf(w, "Welcome to the 'GET /{$}' page !")
	})

	mux.HandleFunc("GET /Login/{$}", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		/* if req.URL.Path != "/LOGIN/" {
			http.NotFound(w, req)
			return
		} */

		c, err := req.Cookie("session_token")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "%v,\nreq.URL.Path= %v\n", c, req.URL.Path)
		fmt.Fprintf(w, "Welcome to the 'GET /Login/{$}' page !")
	})
	log.Fatal(http.ListenAndServe(assets.Port, mux))
}
