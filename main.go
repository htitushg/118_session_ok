package main

import (
	"118_session_ok/assets"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// On relie le fichier css et le favicon au nom static
	fmt.Printf("Main Chemin= %s\n", assets.Chemin+"assets/") //
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(assets.Chemin+"assets/"))))
	http.HandleFunc("/", Home)
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/Signin", Signin)
	http.HandleFunc("/Refresh", Refresh)
	http.HandleFunc("/Logout", Logout)
	http.HandleFunc("/AfficheUserInfo/", AfficheUserInfo)
	// start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
