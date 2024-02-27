package models

import (
	//"Middleware-test/internal/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// In-memory Session data storage
var SessionsData = make(map[string]Session)

func OpenSession(w *http.ResponseWriter, r *http.Request, pseudo string) {

	// Generate and set Session ID cookie
	sessionID := generateSessionID()
	// Generate expiration time for the cookie
	expirationTime := time.Now().Add(time.Minute)

	newCookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expirationTime,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(*w, &newCookie)

	fmt.Printf("%#v\n", newCookie)

	// Create Session data in memory
	SessionsData[sessionID] = Session{
		UserID:         777,
		SessionID:      sessionID,
		Username:       pseudo,
		IpAddress:      GetIP(r),
		ExpirationTime: expirationTime,
	}
}

func RefreshSession(w *http.ResponseWriter, r *http.Request) error {
	// generating new sessionID and new expiration time
	newSessionID := generateSessionID()
	newExpirationTime := time.Now().Add(time.Minute)

	var newCookie = &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		HttpOnly: true,
		Secure:   false, // Use only if using HTTPS
		Path:     "/",
		Expires:  newExpirationTime,
		SameSite: http.SameSiteStrictMode,
	}

	// setting the new cookie
	http.SetCookie(*w, newCookie)

	// retrieving the in-memory current session data
	cookie, err := r.Cookie("session_id")
	currentSessionData := SessionsData[cookie.Value]

	// updating the sessionID and expirationTime
	currentSessionData.SessionID = newSessionID
	currentSessionData.ExpirationTime = newExpirationTime

	// deleting previous entry in the SessionsData map
	delete(SessionsData, cookie.Value)

	// setting the new entry in the SessionsData map
	SessionsData[newSessionID] = currentSessionData

	// adding the new cookie to the request to access it from the targeted handler with the Name "updatedCookie"
	newCookie.Name = "updatedCookie"
	r.AddCookie(newCookie)

	if err != nil {
		return err
	}
	return nil
}

func generateSessionID() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func ValidateSessionID(sessionID string) bool {
	_, ok := SessionsData[sessionID]
	return len(sessionID) == 88 && ok
}
