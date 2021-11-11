package main

import (
	"context"
    "net/http"
)

// takes in our custom handler and converts to http.Handler
func ErrorMiddleware(handler CustomHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			switch e := err.(type) {
			case HTTPError:
				http.Error(w, e.Error(), e.Status())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	})
}

// allows only one type of method to be used on endpoint
func RestrictMethodMiddleware(methods []string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ContainsString(methods, r.Method) {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
		}
        next.ServeHTTP(w, r)
    })
}

// checks auth
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        accessToken := r.Header.Get("x-access-token")

        id, err := ValidateToken(accessToken)
        if err != nil {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        ctx := r.Context()
        ctx = context.WithValue(ctx, "userid", id)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}

