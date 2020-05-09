package main

import (
	"log"
	"net/http"

	"./app"
	"./controller"
	"github.com/gorilla/mux"
)

//  Main function
func main() {
	app.Init()
	r := mux.NewRouter()
	r.HandleFunc("/register", controller.Register).
		Methods("POST")
	r.HandleFunc("/login", controller.Login).
		Methods("POST")
	r.HandleFunc("/profile", controller.ProfileHandler).
		Methods("GET")
	r.HandleFunc("/addmovie", controller.AddMovies).
		Methods("POST")
	r.HandleFunc("/ratemovie", controller.RateCommentMovie).
		Methods("POST")
	r.HandleFunc("/commentmovie", controller.RateCommentMovie).
		Methods("POST")
	r.HandleFunc("/searchmovie", controller.SearchMovie).
		Methods("GET")
	r.HandleFunc("/getmovies", controller.GetMovies).
		Methods("GET")

	log.Fatal(http.ListenAndServe(":9080", r))
}
