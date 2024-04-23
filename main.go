package main

import (
	"fmt"
	"github.com/IDOMATH/session/session"
	"log"
	"net/http"
	"time"
)

var memStore *session.MemoryStore

func main() {
	fmt.Println("Hello world!")
	memStore = session.NewMemoryStore()
	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /session", sessionInsert)
	router.HandleFunc("GET /session", sessionDump)

	log.Fatal(server.ListenAndServe())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home"))
}

func sessionInsert(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("authToken")

	memStore.Insert(token, []byte(token), time.Now().AddDate(0, 0, 1))
	w.Write([]byte("Welcome Home"))
}

//
//func sessionDump(w http.ResponseWriter, r *http.Request) {
//	for token := range memStore.Items {
//		w.Write([]byte(token))
//	}
//}
