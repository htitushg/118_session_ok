package controllers

import (
	"118_session_ok/assets"
	"118_session_ok/models"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Ajouté le 24/02/2024 19h20
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

	//return
}

// Ajout du 24/02/2024 19h20
var logs, _ = os.Create("logs/logs.log")
var jsonHandler = slog.NewJSONHandler(logs, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}).WithAttrs([]slog.Attr{slog.Int("Info", 13)})
var Logger = slog.New(jsonHandler)
var LogId = 0

// Log is a models.Middleware that writes a series of information in logs/logs.log
// in JSON format: time, function name, request Id (incremented int),
// client IP, request Method, and request URL.
func Log() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session_token")
			var pseudo string
			if err != nil {
				pseudo = ""
			} else {
				token := c.Value
				pseudo = assets.Sessions[token].Pseudo
			}
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			//next.ServeHTTP(wrapped, r)

			/* "method", r.Method,
			 */

			LogId++
			log.Println("Log()")
			Logger.Info("Log() Middleware", slog.Int("reqId", LogId), slog.Duration("duration", time.Since(start)), slog.Int("status", wrapped.status), slog.String("path", r.URL.EscapedPath()), slog.String("clientIP", models.GetIP(r)), slog.String("pseudo", pseudo), slog.String("reqMethod", r.Method), slog.String("reqURL", r.URL.String()))
			handler.ServeHTTP(w, r)
		}
	}
}

/* func Log() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fn := func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if err := recover(); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						Logger.Info(
							"err", err,
							"trace", debug.Stack(),
						)
					}
				}()

				LogId++
				log.Println("Log()")
				Logger.Info("Log() Middleware", slog.Int("reqId", LogId), slog.String("clientIP", models.GetIP(r)), slog.String("reqMethod", r.Method), slog.String("reqURL", r.URL.String()))
				handler.ServeHTTP(w, r)

			}
			return fn
		}
	}
} */

// slog.String("Name", GetCurrentName(w, r)),
func Guard() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("Guard()")
			_, exists := SessionValide(w, r)
			if !exists {
				http.Redirect(w, r, "/Login?err=restricted", http.StatusSeeOther)
				return
			}
			handler.ServeHTTP(w, r)
		}
	}
}

func Foo() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("Foo()")
			handler.ServeHTTP(w, r)
		}
	}
}

// Join is used to concatenate various middlewares, for better visibility.
// it takes the http.HandlerFunc corresponding to the route, and then
// any number of models.Middleware that will be concatenated in order like this:
// middlewares[0](middlewares[1](middlewares[2](handlerFunc))).
func Join(handlerFunc http.HandlerFunc, middlewares ...models.Middleware) http.HandlerFunc {
	if len(middlewares) == 1 {
		return middlewares[0](handlerFunc)
	}
	return middlewares[0](Join(handlerFunc, middlewares[1:]...))
}
