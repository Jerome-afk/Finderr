package main

import (
	routes "finderr/routeHandler"
	"fmt"
	"log"
	"net/http"
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
    router := routes.InitServer()
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