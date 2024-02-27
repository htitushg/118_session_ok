package main

import (
	"118_session_ok/assets"
	"118_session_ok/controllers"
	"fmt"
	"os"
	"time"

	"log"
	"net/http"
)

func main() {
	// If the file doesn't exist, create it or append to the file
	location, _ := time.LoadLocation("France/Paris")
	fmt.Println(location)
	file, err := os.Create(assets.Chemin + "logs/logs.txt") //, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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
	mux.HandleFunc("GET /Login", controllers.LoginBundle)
	mux.HandleFunc("POST /Signin", controllers.SigninBundle)
	mux.HandleFunc("GET /Index", controllers.IndexBundle)
	mux.HandleFunc("POST /Logout", controllers.LogoutBundle)
	mux.HandleFunc("GET /Register", controllers.RegisterBundle)
	mux.HandleFunc("POST /AfficheUserInfo", controllers.AfficheUserInfoBundle)

	// Handling MethodNotAllowed error on /
	mux.HandleFunc("/{$}", controllers.IndexHandlerNoMethBundle)

	// Handling StatusNotFound error everywhere else
	mux.HandleFunc("/", controllers.IndexHandlerOtherBundle)
	// start the server
	fmt.Printf("http://localhost%v , ctrl+clic sur le lien pour lancer le navigateur\n", assets.Port)
	log.Fatal(http.ListenAndServe(assets.Port, mux))
}
