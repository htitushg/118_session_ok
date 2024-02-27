package controllers

import (
	//"Middleware-test/models"
	"118_session_ok/models"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

var logs, _ = os.Create("logs/logs.log")
var jsonHandler = slog.NewJSONHandler(logs, nil)
var Logger = slog.New(jsonHandler)
var LogId = 0

// Log is a models.Middleware that writes a series of information in logs/logs.log
// in JSON format: time, function name, request Id (incremented int),
// client IP, request Method, and request URL.
var Log models.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_id")
		var pseudo string
		if err != nil {
			pseudo = ""
		} else {
			token := c.Value
			pseudo = models.SessionsData[token].Username
		}
		start := time.Now()
		wrapped := wrapResponseWriter(w)
		LogId++
		log.Println("Log()")
		Logger.Info("Log() Middleware", slog.Int("reqId", LogId), slog.Duration("duration", time.Since(start)), slog.Int("status", wrapped.status), slog.String("path", r.URL.EscapedPath()), slog.String("clientIP", models.GetIP(r)), slog.String("pseudo", pseudo), slog.String("reqMethod", r.Method), slog.String("reqURL", r.URL.String()))
		next.ServeHTTP(w, r)
	}
}

// Guard is a models.Middleware that verify if a user has an opened session
// through the cookies and let it pass if ok, and redirects if not.
var Guard models.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Guard()")
		// Extract session ID from cookie
		cookie, err := r.Cookie("session_id")
		if err != nil || !models.ValidateSessionID(cookie.Value) {
			// Handle invalid session (e.g., redirect to login)
			Logger.Warn("Invalid session", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusUnauthorized))
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Retrieve user data from session
		userData, ok := models.SessionsData[cookie.Value]
		if !ok {
			// Handle missing session (e.g., redirect to login)
			Logger.Warn("Invalid session", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusUnauthorized))
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Verify user IP address
		if userData.IpAddress != models.GetIP(r) {
			// Handle missing session (e.g., redirect to login)
			Logger.Warn("Invalid session", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusUnauthorized))
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Verify expiration time
		fmt.Printf("%#v\n", cookie)
		if userData.ExpirationTime.Before(time.Now()) {
			// Handle missing session (e.g., redirect to login)
			Logger.Warn("Invalid session", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusUnauthorized))
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		err = models.RefreshSession(&w, r)
		if err != nil {
			Logger.Error(err.Error())
		}
		var cookies []http.Cookie
		cookie1, _ := r.Cookie("session_id")
		cookie2, _ := r.Cookie("updatedCookie")
		cookies = append(cookies, *cookie1, *cookie2)
		log.Printf("Guard() Middleware, after RefreshSession(): %#v\n", cookies)
		// Use user data (e.g., display username)
		//fmt.Fprintf(w, "Welcome, user %s", userData["user_id"])
		next.ServeHTTP(w, r)
	}
}

// Foo is a random models.Middleware for tests
var Foo models.Middleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Foo()")
		next.ServeHTTP(w, r)
	}
}

// Join is used to concatenate various middlewares, for better visibility.
// it takes the http.HandlerFunc corresponding to the route, and then
// any number of models.Middleware that will be concatenated in order like this:
// middlewares[0](middlewares[1](middlewares[2](handlerFunc))).
func Join(handlerFunc http.HandlerFunc, middlewares ...models.Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handlerFunc = middlewares[i](handlerFunc)
	}
	return handlerFunc
}
