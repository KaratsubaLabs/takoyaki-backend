package main

import (
    _ "fmt"
    "net/http"
)

func main() {

    port := ":8080"

    var mux *http.ServeMux = http.NewServeMux()
    Routes(mux)

    fmt.Printf("Listening on port %s\n", port)
    http.ListenAndServe(port, handler)

}

