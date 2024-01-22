package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"118_session/assets"
	"118_session/data"
	"118_session/helpers"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Handlers
// log.Printf("log: UrlPath: %#v\n", r.URL.Path) // testing
//
//	if r.URL.Path != "/create/treatment" {
//		errorHandler(w, r, http.StatusNotFound)
//		return
//	}
//
// for GET
func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("Home log: UrlPath: %#v\n", r.URL.Path) // testing
	if data.Sessionvalid(w, r) {
		http.Redirect(w, r, "/Index", http.StatusFound)
	} else {
		fmt.Printf("Home Chemin= %s\n", assets.Chemin+"templates/home.html")
		var body, _ = helpers.LoadFile(assets.Chemin + "templates/home.html")
		fmt.Fprint(w, body)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Login log: UrlPath: %#v\n", r.URL.Path) // testing
	if data.Sessionvalid(w, r) {
		c, err := r.Cookie("session_token")
		assets.CheckError(err)
		sessionToken := c.Value
		// remove the users session from the session map
		delete(assets.Sessions, sessionToken)

		// We need to let the client know that the cookie is expired
		// In the w, we set the session token to an empty
		// value and set its expiry as the current time
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		fmt.Printf("Login Chemin= %s\n", assets.Chemin+"templates/login.html")
		var body, _ = helpers.LoadFile(assets.Chemin + "templates/login.html")
		fmt.Fprint(w, body)
	}
}

// for POST
func LoginPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginPost log: UrlPath: %#v\n", r.URL.Path) // testing
	fmt.Println("LoginPost")
	creds := &assets.CredentialsR{
		Pseudo:   r.FormValue("pseudo"),
		Password: r.FormValue("passid"),
	}
	fmt.Printf("Name= %s, Password= %s\n", creds.Pseudo, creds.Password)
	redirectTarget := "/"
	if len(creds.Pseudo) > 0 && len(creds.Password) > 0 {
		// Database check for user data!
		_userIsValid := data.UserIsValid(*creds)

		if _userIsValid {
			// Create a new random session token
			// we use the "github.com/google/uuid" library to generate UUIDs
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(120 * time.Second)

			// Set the token in the session map, along with the session information
			assets.Sessions[sessionToken] = assets.Session{
				Pseudo: creds.Pseudo,
				Expiry: expiresAt,
			}

			// Finally, we set the client cookie for "session_token" as the session token we just generated
			// we also set an expiry time of 120 seconds
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: expiresAt,
			})
			redirectTarget = "/Index"

		} else {
			redirectTarget = "/Register"
		}
	}
	http.Redirect(w, r, redirectTarget, http.StatusFound)
}

// for GET
func Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register log: UrlPath: %#v\n", r.URL.Path) // testing
	if data.Sessionvalid(w, r) {
		http.Redirect(w, r, "/Index", http.StatusFound)
	} else {
		templates, err := template.ParseFiles(assets.Chemin + "templates/register3.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var rcreds assets.Credentials
		if err := templates.Execute(w, rcreds); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// for POST
func RegisterPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("RegisterPost log: UrlPath: %#v\n", r.URL.Path) // testing
	fmt.Println("RegisterPost")
	var rpcreds assets.Credentials
	r.ParseForm()
	rpcreds.Pseudo = r.FormValue("pseudo")
	rpcreds.Email = r.FormValue("email")
	rpcreds.Password = r.FormValue("passid")
	rpcreds.Password2 = r.FormValue("passid2")
	rpcreds.Firstname = r.FormValue("firstname")
	rpcreds.Lastname = r.FormValue("lastname")
	rpcreds.Address = r.FormValue("address")
	rpcreds.Town = r.FormValue("town")
	rpcreds.ZipCode = r.FormValue("zip")
	rpcreds.Country = r.FormValue("country")
	rpcreds.Genre = r.FormValue("sex")
	rpcreds.Description = r.FormValue("desc")
	rpcreds.Message = ""

	if rpcreds.Password == rpcreds.Password2 {
		fmt.Printf("pseudo = %s, password= %s, confirmpassword= %s\n", rpcreds.Pseudo, rpcreds.Password, rpcreds.Password2)
		_uName, _pwd, _email := false, false, false
		_uName = !helpers.IsEmpty(rpcreds.Pseudo)
		_pwd = !helpers.IsEmpty(rpcreds.Password)
		_email = !helpers.IsEmpty(rpcreds.Email)
		if _uName && _pwd && _email {
			isCreate := data.UserCreate(rpcreds)
			if isCreate {
				var rpcredsR assets.CredentialsR
				rpcredsR.Pseudo = rpcreds.Pseudo
				rpcredsR.Email = rpcreds.Email
				rpcredsR.Password = rpcreds.Password
				rpcredsR.Firstname = rpcreds.Firstname
				rpcredsR.Lastname = rpcreds.Lastname

				fmt.Printf("RegisterPost Chemin = %s\n", assets.Chemin+"templates/createuser.html")
				t, err := template.ParseFiles(assets.Chemin + "templates/createuser.html")
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				if err := t.Execute(w, rpcredsR); err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			} else {
				rpcreds.Message = "Il n'a pas été possible de créer l'utilisateur ou l'utilisateur ou l'adresse mail existe déjà!"
				t, err := template.ParseFiles(assets.Chemin + "templates/register3.html")
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				if err := t.Execute(w, rpcreds); err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
		} else {
			//fmt.Fprintln(w, "This fields can not be blank!")
			rpcreds.Message = "This fields can not be blank!"
			var body, _ = helpers.LoadFile(assets.Chemin + "templates/register3.html")
			fmt.Fprint(w, body)
		}

	} else {
		//fmt.Fprintln(w, "Les mots de passe doivent être identiques")
		rpcreds.Message = "Les mots de passe doivent être identiques"
		var body, _ = helpers.LoadFile(assets.Chemin + "templates/register3.html")
		fmt.Fprint(w, body)
	}
}

// for GET
func Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("Index log: UrlPath: %#v\n", r.URL.Path) // testing
	fmt.Println("Index")
	redirectTarget := "templates/home.html"
	if !data.Sessionvalid(w, r) {
		redirectTarget = "templates/home.html"
	} else {

		//fmt.Println("Expiry=", assets.Session.Expiry)

		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				redirectTarget = "templates/home.html"
				//w.WriteHeader(http.StatusUnauthorized)
				//return
			} else {
				// For any other type of error, return a bad r status
				redirectTarget = "templates/home.html"
				//w.WriteHeader(http.StatusBadr)
			}
		} else {
			sessionToken := c.Value
			//assets.Sessions[sessionToken]

			if !helpers.IsEmpty(assets.Sessions[sessionToken].Pseudo) {
				if assets.Sessions[sessionToken].Expiry.Before(time.Now()) {
					redirectTarget = "/home.html"
				} else {
					DatedeCreation := assets.Sessions[sessionToken].Expiry.Format("“2006-01-02 15h04 05 secondes”")
					DJour := time.Now().Format("2006-01-02")
					type data struct {
						CSessions   assets.Session
						Date_Expire string
						Date_jour   string
						SToken      string
					}
					Data := data{
						CSessions:   assets.Sessions[sessionToken],
						Date_Expire: DatedeCreation,
						Date_jour:   DJour,
						SToken:      sessionToken,
					}
					fmt.Printf("IndexPageHandler 114 Chemin = %s\n", assets.Chemin+"templates/index.html")
					t, err := template.ParseFiles(assets.Chemin + "templates/index.html")
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					err = t.Execute(w, Data)
					if err != nil {
						http.Error(w, err.Error(), 500)
						return
					}
					redirectTarget = "templates/Index"
				}
			} else {
				redirectTarget = "templates/home.html"
			}
		}
	}
	fmt.Printf("redirectTarget: %v\n", redirectTarget)
	//http.Redirect(w, r, redirectTarget, http.StatusFound)
	var body, _ = helpers.LoadFile(assets.Chemin + redirectTarget)
	fmt.Fprint(w, body)
}
func Deconnexion(w http.ResponseWriter, r *http.Request) {
	log.Printf("Deconnexion log: UrlPath: %#v\n", r.URL.Path) // testing
	fmt.Println("Deconnexion")
	nom := strings.TrimPrefix(r.URL.Path, "/Deconnexion/")
	fmt.Println("Nom: ", nom)
	fmt.Printf("Deconnexion Chemin = %s\n", assets.Chemin+"templates/deconnexion.html")
	t, err := template.ParseFiles(assets.Chemin + "templates/deconnexion.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, nom); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func DeconnexionPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeconnexionPost log: UrlPath: %#v\n", r.URL.Path) // testing
	if data.Sessionvalid(w, r) {
		c, err := r.Cookie("session_token")
		assets.CheckError(err)
		sessionToken := c.Value
		// remove the users session from the session map
		delete(assets.Sessions, sessionToken)

		// We need to let the client know that the cookie is expired
		// In the w, we set the session token to an empty
		// value and set its expiry as the current time
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// Cookie

func SetCookie(userName string, w http.ResponseWriter) {
	fmt.Println("SetCookie")
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func ClearCookie(w http.ResponseWriter) {
	fmt.Println("ClearCookie")
	cookie := &http.Cookie{
		Name:   "cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func GetUserName(r *http.Request) (userName string) {
	fmt.Println("GetUserName")
	if cookie, err := r.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}
