package main

import (
    "fmt"
    "net/http"
)

const PORT = "8080"
func StartServer() {

    var mux *http.ServeMux = http.NewServeMux()
    Routes(mux)

    fmt.Printf("Listening on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, mux)

}

