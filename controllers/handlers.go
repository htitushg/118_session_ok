package controllers

import (
	"118_session_ok/assets"
	"118_session_ok/models"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// Controlleur Home: Affiche le Page publique(home) si la session n'est pas valide, sinon affiche la page privée(index)
func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("Home log: UrlPath: %#v\n", r.URL.Path) // testing
	log.Printf("%#v\n", models.GetCurrentFuncName())
	var message template.HTML
	switch r.URL.Query().Get("err") {
	case "pass":
		message = "<div class=\"message\">Wrong username or password!</div>"
	case "restricted":
		message = "<div class=\"message\">You need to sign in to access to this resource!</div>"
	}
	tmpl, err := template.ParseFiles(assets.Chemin + "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err = tmpl.ExecuteTemplate(w, "home", message); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Login: Affiche la page de connexion
func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login log: UrlPath: %#v\n", r.URL.Path)
	log.Printf("%#v\n", models.GetCurrentFuncName())
	var message template.HTML
	switch r.URL.Query().Get("err") {
	case "pass":
		message = "<div class=\"message\">Wrong username or password!</div>"
	case "restricted":
		message = "<div class=\"message\">You need to sign in to access to this resource!</div>"
	}
	tmpl, err := template.ParseFiles(assets.Chemin + "templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "login", message); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Signin: Traite les informations de connexin fournies par login,
// Si les informations sont incorrectes, renvoie vers home
// Sinon crée la session et renvoie vers index
func Signin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Signin log: UrlPath: %#v\n", r.URL.Path)
	log.Printf("%#v\n", models.GetCurrentFuncName())
	var creds assets.Credentials
	creds.Pseudo = r.FormValue("pseudo")
	creds.Password = r.FormValue("passid")
	// Get the expected password from our in memory map
	expectedPassword, ok := assets.Users[creds.Pseudo]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/Login?err=pass", http.StatusSeeOther)
		return
	}
	models.OpenSession(&w, r, creds.Pseudo)
	http.Redirect(w, r, "/Index", http.StatusSeeOther)
}
func Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("Index log: UrlPath: %#v\n", r.URL.Path)
	log.Printf("%#v\n", models.GetCurrentFuncName())
	sessionID, _ := r.Cookie("updatedCookie")
	var data = assets.Data{
		SToken:    sessionID.Value,
		CSessions: models.SessionsData[sessionID.Value],
		Date_jour: time.Now().Format("2006-01-02"),
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Logout: Si la session est valide, Ferme la session et renvoie vers home
func Logout(w http.ResponseWriter, r *http.Request) {
	log.Printf("Logout log: UrlPath: %#v\n", r.URL.Path)
	log.Printf("%#v\n", models.GetCurrentFuncName())
	var err error
	c, err := r.Cookie("session_token")
	if err == nil {
		sessionToken := c.Value

		// remove the users session from the session map
		delete(models.SessionsData, sessionToken)
		// We need to let the client know that the cookie is expired
		// In the response, we set the session token to an empty
		// value and set its expiry as the current time
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
		})
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.ExecuteTemplate(w, "home", nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur AfficheUserInfo: Si la session est valide renvoie vers afficheuserinfo
// Sinon renvoie vers home
func AfficheUserInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("AfficheUserInfo log: UrlPath: %#v\n", r.URL.Path) // testing
	log.Println(models.GetCurrentFuncName())
	cookie, err := r.Cookie("updatedCookie")
	if err != nil {
		http.Redirect(w, r, "/Home?err=pass", http.StatusSeeOther)
	}
	var sessionToken = cookie.Value
	var data = assets.Data{
		SToken:      sessionToken,
		CSessions:   models.SessionsData[sessionToken],
		Date_jour:   time.Now().Format("2006-01-02"),
		Date_Expire: models.SessionsData[sessionToken].ExpirationTime,
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/afficheuserinfo.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.ExecuteTemplate(w, "afficheuserinfo", data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Register: Renvoie vers register pour enregistrement
func Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register log: UrlPath: %#v\n", r.URL.Path) // testing
	log.Println(models.GetCurrentFuncName())
	//var data assets.Data
	var t *template.Template
	_, err := r.Cookie("updatedCookie")
	if err != nil {
		t, err = template.ParseFiles(assets.Chemin + "templates/register.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := t.ExecuteTemplate(w, "register", nil); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Redirect(w, r, "/Home?err=pass", http.StatusSeeOther)
	}
}
func IndexHandlerNoMeth(w http.ResponseWriter, r *http.Request) {
	log.Println(models.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusMethodNotAllowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	Logger.Warn("indexHandlerNoMeth", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusMethodNotAllowed))
	http.Redirect(w, r, "/Home?err=pass", http.StatusSeeOther)
	//http.Redirect(w, r, "/Home", http.StatusSeeOther)
}
func IndexHandlerOther(w http.ResponseWriter, r *http.Request) {
	log.Println(models.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	Logger.Warn("indexHandlerOther", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusNotFound))
	http.Redirect(w, r, "/Home?err=pass", http.StatusSeeOther)
	//http.Redirect(w, r, "/Home", http.StatusSeeOther)
}
