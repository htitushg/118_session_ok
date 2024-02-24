package main

import (
	"118_session_ok/assets"
	"118_session_ok/controllers"
	"fmt"
	"os"

	"log"
	"net/http"
)

func main() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(assets.Chemin+"logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	mux := http.NewServeMux()
	// On relie le fichier css et le favicon au nom static
	log.Printf("Main Chemin= %s\n", assets.Chemin+"assets/") //
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(assets.Chemin+"assets/"))))
	// routes

	mux.HandleFunc("GET /{$}", controllers.HomeBundle)
	mux.HandleFunc("POST /Login", controllers.LoginBundle)
	mux.HandleFunc("POST /Signin", controllers.SigninBundle)
	mux.HandleFunc("POST /Logout", controllers.LogoutBundle)
	mux.HandleFunc("POST /Register", controllers.RegisterBundle)
	mux.HandleFunc("POST /AfficheUserInfo", controllers.LogMiddleware(controllers.AfficheUserInfoBundle))

	// Handling MethodNotAllowed error on /
	mux.HandleFunc("/{$}", controllers.IndexHandlerNoMethBundle)

	// Handling StatusNotFound error everywhere else
	mux.HandleFunc("/", controllers.IndexHandlerOtherBundle)
	// start the server
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur", assets.Port)
	log.Fatal(http.ListenAndServe(assets.Port, mux))
}
