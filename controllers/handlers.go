package controllers

import (
	"118_session_ok/assets"
	"118_session_ok/models"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// Ajouté le 02/02/2024
// Note - NOT RFC4122 compliant

// Génère un UUID (Jeton de session : Token)
func Pseudo_uuid() (uuid string) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	/* uuid = hex.EncodeToString(b)
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X-%X", uuid[0:10], uuid[10:20], uuid[20:30], uuid[30:40], uuid[40:50], uuid[50:]) */
	uuid = fmt.Sprintf("%X-%X-%X-%X", b[0:10], b[10:20], b[20:30], b[30:])
	return
}

// Si la session est valide, renvoie le Token et true, sinon nil et false
func SessionExiste(w http.ResponseWriter, r *http.Request) (stoken string, resultat bool) {
	log.Printf("SessionExiste r.Method= %v\n", r.Method)
	resultat = false
	stoken = ""
	c, err := r.Cookie("session_token")
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
	_, exists := models.SessionsData[stoken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	}
	resultat = true
	return stoken, resultat
}
func SessionValide(w http.ResponseWriter, r *http.Request) (stoken string, resultat bool) {
	log.Printf("SessionValide r.Method= %v\n", r.Method)
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
	_, exists := models.SessionsData[stoken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	}
	// If the previous session is valid, create a new session token for the current user
	// on peut utiliser google : "github.com/google/uuid"
	// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

	/* newSessionToken := uuid.NewString() */
	/* newSessionToken := Pseudo_uuid()
	maxAge := 120 */

	// Set the token in the session map, along with the user whom it represents
	/* assets.Sessions[newSessionToken] = assets.Session{
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
	}) */
	resultat = true
	return stoken, resultat
}

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
	//sessionToken, exists := SessionExiste(w, r)
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
