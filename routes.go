package main

import (
    "net/http"
	"encoding/json"
)

type routeInfo struct {
	route        string
	methods      []string // possibly restrict to certain strings (ie POST, GET)
	authRoute    bool
	handlerFn    func(http.ResponseWriter, *http.Request) (int, error)
}

func (info routeInfo) routeHandler(w http.ResponseWriter, r *http.Request) {

	statusCode, err := func() (int, error) {

		// restrict method of request
		err := _
		if err != nil {
			return http.MethodNotAllowed, err
		}

		// restrict auth
		err := _
		if err != nil {
			return http.StatusUnauthorized, err
		}

		// validate request (if possible)
		err := _
		if err != nil {
			return http.StatusBadRequest, err
		}

		// delegate to handler
		statusCode, err := info.handlerFn(w, r)
		if err != nil {
			return statusCode, err
		}
	}()

	// take care of errors / return code + log
	switch statusCode {
	default:
	}

}

var routeSchema = []routeInfo{
	{
		route: "/register",
		methods: []string{"POST"},
		authRoute: false,
		handlerFn: registerHandler,
	},
}

type registerRequest struct {
	Username      string
	Password      string
	Email         string
}
func registerHandler(w http.ResponseWriter, r *http.Request) (int, error) {

	// (possibly have db connection be part of the context)
	db, err := DBConnection()
	if err != nil {
		return http.StatusInternalServerError, "could not connect to database"
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
	for _, routeInfo : range routeSchema {
		mux.Handle(routeInfo.route, routeInfo)
	}
}

