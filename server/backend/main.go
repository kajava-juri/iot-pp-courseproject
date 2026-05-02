package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    // 1. Initialize the Chi router
    r := chi.NewRouter()

    // 2. Add Middleware (optional but recommended)
    r.Use(middleware.Logger)

    // 3. Create a Dummy API (GET Endpoint)
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })

    // 4. Run the Server
    http.ListenAndServe(":8080", r)
}