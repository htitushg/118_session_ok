package main

import (
	"118_session_ok/assets"
	"crypto/rand"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Ajouté le 02/02/2024
// Note - NOT RFC4122 compliant
func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

// Fin de l'Ajout du 02/02/2024
func SessionValide(w http.ResponseWriter, r *http.Request) (stoken string, resultat bool) {
	c, err := r.Cookie("session_token")
	resultat = false
	stoken = ""
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return stoken, resultat
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return stoken, resultat
	}
	stoken = c.Value
	_, exists := assets.Sessions[stoken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	}
	// If the previous session is valid, create a new session token for the current user
	// on peut utiliser google : "github.com/google/uuid"
	// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

	/* newSessionToken := uuid.NewString() */
	newSessionToken := pseudo_uuid()
	maxAge := 120

	// Set the token in the session map, along with the user whom it represents
	assets.Sessions[newSessionToken] = assets.Session{
		Pseudo: assets.Sessions[stoken].Pseudo,
		MaxAge: maxAge,
	}
	// Delete the older session token
	delete(assets.Sessions, stoken)
	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  newSessionToken,
		MaxAge: maxAge,
	})
	/* if assets.Sessions[stoken].Expiry.Before(time.Now()) {
		delete(assets.Sessions, stoken)
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	} */
	resultat = true
	return newSessionToken, resultat
}
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Home log: UrlPath: %#v\n", r.URL.Path) // testing
	var data assets.Data
	var t *template.Template
	var err error
	stoken, exists := SessionValide(w, r)
	if !exists {
		t, err = template.ParseFiles(assets.Chemin + "templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		// Il nous faut ici rassembler les infos utilisateur
		DJour := time.Now().Format("2006-01-02")
		data.CSessions = assets.Sessions[stoken]
		data.Date_jour = DJour
		data.SToken = stoken
		/* data.Email:       credsR.Email
		data.Firstname:   credsR.Firstname
		data.Lastname:   credsR.Lastname */

		t, err = template.ParseFiles(assets.Chemin + "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err = t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Login log: UrlPath: %#v\n", r.URL.Path)
	t, err := template.ParseFiles(assets.Chemin + "templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds assets.Credentials
	var data assets.Data
	var t *template.Template
	var err error
	fmt.Printf("Signin log: UrlPath: %#v\n", r.URL.Path)
	creds.Pseudo = r.FormValue("pseudo")
	creds.Password = r.FormValue("passid")
	// Get the expected password from our in memory map
	expectedPassword, ok := assets.Users[creds.Pseudo]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		t, err = template.ParseFiles(assets.Chemin + "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

	} else {

		// Create a new random session token
		// on peut utiliser google : "github.com/google/uuid"
		// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

		/* newSessionToken := uuid.NewString() */
		sessionToken := pseudo_uuid()
		maxAge := 120

		// Set the token in the session map, along with the user whom it represents
		assets.Sessions[sessionToken] = assets.Session{
			Pseudo: creds.Pseudo,
			MaxAge: maxAge,
		}

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  sessionToken,
			MaxAge: maxAge,
		})
		//DatedeCreation := assets.Sessions[sessionToken].Expiry.Format("“2006-01-02 15h04 05 secondes”")
		DJour := time.Now().Format("2006-01-02")
		data.CSessions = assets.Sessions[sessionToken]
		data.Date_jour = DJour
		data.SToken = sessionToken

		t, err = template.ParseFiles(assets.Chemin + "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionToken, exists := SessionValide(w, r)
	if exists {
		// remove the users session from the session map
		delete(assets.Sessions, sessionToken)
		// We need to let the client know that the cookie is expired
		// In the response, we set the session token to an empty
		// value and set its expiry as the current time
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func AfficheUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AfficheUserInfo log: UrlPath: %#v\n", r.URL.Path) // testing
	var data assets.Data
	var err error
	var t *template.Template
	sessionToken, exists := SessionValide(w, r)
	if exists {
		DJour := time.Now().Format("2006-01-02")
		data.CSessions = assets.Sessions[sessionToken]
		data.Date_jour = DJour
		data.SToken = sessionToken
		t, err = template.ParseFiles(assets.Chemin + "templates/afficheuserinfo.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		t, err = template.ParseFiles(assets.Chemin + "templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Ajouté le 02/02/2024
// for GET
func Register(w http.ResponseWriter, r *http.Request) {
	var err error
	var t *template.Template
	fmt.Printf("Register log: UrlPath: %#v\n", r.URL.Path) // testing
	t, err = template.ParseFiles(assets.Chemin + "templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
