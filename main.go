package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello world!")
	router := http.NewServeMux()

	router.HandleFunc("GET /", handleHome)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home"))
}
