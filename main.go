package main

import (
	"fmt"
	"net/http"
	"os"

	"finderr/functions"
)

func main() {
	myargs := os.Args[1:]
	if len(myargs) > 0 {
		fmt.Println("Too many args passed")
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/",  functions.Router)
}
