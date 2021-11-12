package main

import (
    "net/http"
	_ "encoding/json"
)

type CustomHandler = func(http.ResponseWriter, *http.Request) error

type routeInfo struct {
	route        string
	methods      []string // possibly restrict to certain strings (ie POST, GET)
	authRoute    bool
	bodySchema   interface{}
	handlerFn    CustomHandler
}

func (info routeInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var handlerWithMiddleware http.Handler = ErrorMiddleware(info.handlerFn)

    // parse body (if applicable)
    if info.bodySchema != nil {
        handlerWithMiddleware = ParseBodyJSONMiddleware(info.bodySchema, handlerWithMiddleware)
    }

	// restrict auth (if applicable)
	if info.authRoute {
		handlerWithMiddleware = AuthMiddleware(handlerWithMiddleware)
	}

	// restrict method of request
	handlerWithMiddleware = RestrictMethodMiddleware(info.methods, handlerWithMiddleware)

	// delegate to handler
	handlerWithMiddleware.ServeHTTP(w, r)

}

var routeSchema = []routeInfo{
	{
		route: "/register",
		methods: []string{"POST"},
		authRoute: false,
		bodySchema: registerRequest{},
		handlerFn: registerHandler,
	},
}

type registerRequest struct {
	Username      string
	Password      string
	Email         string
}
func registerHandler(w http.ResponseWriter, r *http.Request) error {

	// (possibly have db connection be part of the context)
	db, err := DBConnection()
	if err != nil {
        return HTTPStatusError{http.StatusInternalServerError, err}
	}
	defer db.Close()

    return nil
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
	for _, routeInfo := range routeSchema {
		mux.Handle(routeInfo.route, routeInfo)
	}
}

