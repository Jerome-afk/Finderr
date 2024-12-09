package main

import (
	"finderr/server"
	"log"
	"net/http"
    "fmt"
	"os"
	"github.com/joho/godotenv"
)

 func init()  {
     err := godotenv.Load(".env")
     if err != nil {
         fmt.Println("Error loading .env file")
         
    }
 }

func main() {
    // Initialize router
    router := server.InitServer()
    port := os.Getenv("PORT")
    if port == "" {
        port = "8787"
    }

    log.Printf("Server starting on port %s...", port)
    err := http.ListenAndServe(":"+ port, router)
    if err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}