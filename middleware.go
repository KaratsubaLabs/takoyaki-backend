package main

import (
	"fmt"
	"context"
    "net/http"
)

// allows only one type of method to be used on endpoint
func RestrictMethod(method string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != method {
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func RestrictMethodGet(next http.Handler) http.Handler {
	return RestrictMethod("GET", next)
}

func RestrictMethodPost(next http.Handler) http.Handler {
	return RestrictMethod("POST", next)
}

// checks auth
func RestrictAuth(next http.Handler) http.Handler {
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

