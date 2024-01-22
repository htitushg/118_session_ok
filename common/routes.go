package handlers

import (
	"net/http"
)

func Routes() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/Index", Index)
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/LoginPost", LoginPost)
	http.HandleFunc("/Register", Register)
	http.HandleFunc("/RegisterPost", RegisterPost)
	http.HandleFunc("/Deconnexion/", Deconnexion)
	http.HandleFunc("/DeconnexionPost", DeconnexionPost)
}
