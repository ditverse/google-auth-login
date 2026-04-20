package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google-auth-backend/handler"
	"google-auth-backend/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/google", middleware.CORS(handler.GoogleLogin))
	mux.HandleFunc("/auth/me", middleware.CORS(handler.GetMe))
	mux.HandleFunc("/auth/logout", middleware.CORS(handler.Logout))

	fmt.Printf("Server berjalan di http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
