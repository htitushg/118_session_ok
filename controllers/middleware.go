package controllers

import (
	"118_session_ok/models"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var logs, _ = os.Create("logs/logs.log")
var jsonHandler = slog.NewJSONHandler(logs, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}).WithAttrs([]slog.Attr{
	slog.Int("Info", 13),
})
var Logger = slog.New(jsonHandler)
var LogId = 0

// Log is a models.Middleware that writes a series of information in logs/logs.log
// in JSON format: time, function name, request Id (incremented int),
// client IP, request Method, and request URL.
func Log() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			LogId++
			log.Println("Log()")
			Logger.Info("Log() Middleware", slog.Int("reqId", LogId), slog.String("clientIP", models.GetIP(r)), slog.String("reqMethod", r.Method), slog.String("reqURL", r.URL.String()))
			handler.ServeHTTP(w, r)
		}
	}
}

func Guard() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("Guard()")
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
