package main

import (
	"118_session_ok/assets"
	"fmt"

	"log"
	"net/http"
	"os"
)

func main() {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	mux := http.NewServeMux()
	// On relie le fichier css et le favicon au nom static
	log.Printf("Main Chemin= %s\n", assets.Chemin+"assets/") //
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(assets.Chemin+"assets/"))))
	// routes

	mux.HandleFunc("GET /{$}", LogMiddleware(Home))
	mux.HandleFunc("POST /Login", LogMiddleware(Login))
	mux.HandleFunc("POST /Signin", LogMiddleware(Signin))
	mux.HandleFunc("POST /Logout", LogMiddleware(Logout))
	mux.HandleFunc("POST /Register", LogMiddleware(Register))
	mux.HandleFunc("POST /AfficheUserInfo", LogMiddleware(AfficheUserInfo))
	// start the server
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur", assets.Port)
	log.Fatal(http.ListenAndServe(assets.Port, mux))
}
