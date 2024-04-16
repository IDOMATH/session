package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var memStore = NewMemoryStore()

func main() {
	fmt.Println("Hello world!")
	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.HandleFunc("GET /", handleHome)

	log.Fatal(server.ListenAndServe())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("authToken")

	memStore.Insert(strconv.Itoa(rand.Int()), []byte(token), time.Now().AddDate(0, 0, 1))
	w.Write([]byte("Welcome Home"))
}
