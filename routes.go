package main

import (
    "net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
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
    mux.Handle("/register", http.HandlerFunc(registerHandler))
    mux.Handle("/login", http.HandlerFunc(loginHandler))
}
