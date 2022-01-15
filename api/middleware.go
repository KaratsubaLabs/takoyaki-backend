package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
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

// validate body (must be used after ParseBodyJsonMiddleware
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		parsedBody, ok := r.Context().Value(ContextKeyParsedBody).(interface{})
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// could make validator a global object (might be a bit more performant)
		err := validator.New().Struct(parsedBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")

			validationErrors := err.(validator.ValidationErrors)
			json.NewEncoder(w).Encode(validationErrors)
			// fmt.Printf("Validation Errors %+v", validationErrors)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// parses body as json
func ParseBodyJSONMiddleware(bodySchema interface{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		parsedBody := bodySchema
		err := dec.Decode(&parsedBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyParsedBody, parsedBody)
		r = r.WithContext(ctx)

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
		ctx = context.WithValue(ctx, ContextKeyUserID, id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[%s] %s\n", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		next.ServeHTTP(w, r)
	})
}
