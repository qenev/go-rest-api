package main

import (
    "fmt"
    "log"
    "net/http"

    "go-rest-api/handlers"
)

func main() {
    store := handlers.NewStore()
    mux   := http.NewServeMux()

    mux.HandleFunc("/api/items",    store.HandleItems)
    mux.HandleFunc("/api/items/",   store.HandleItem)

    addr := ":8080"
    fmt.Printf("Server running on http://localhost%s\n", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}
