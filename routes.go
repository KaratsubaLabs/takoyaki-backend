package main

import (
    "net/http"
)

type registerRequest struct {
	Username      string
	Password      string
	Email         string
}
func registerHandler(w http.ResponseWriter, r *http.Request) {

	// validation first

	db, err := DBConnection()
	if err != nil {
		http.Error(w, "could not connect to database", http.StatusInternalServerError)
	}
	defer db.Close()



}

type loginRequest struct {
	Username      string
	Password      string
	Email         string
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
}

func infoVPSHandler(w http.ResponseWriter, r *http.Request) {
}

func createVPSHandler(w http.ResponseWriter, r *http.Request) {
}

func destroyVPSHandler(w http.ResponseWriter, r *http.Request) {
}

func Routes(mux *http.ServeMux) {
    mux.Handle("/register",
		RestrictMethodPost(http.HandlerFunc(registerHandler))
	)
    mux.Handle("/login",
		RestrictMethodPost(http.HandlerFunc(loginHandler))
	)
}

