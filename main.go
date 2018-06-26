package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cmcgow061408/book-microservice/api"
)

func main() {
	fmt.Println("Starting service....")
	http.HandleFunc("/", index)
	http.HandleFunc("/api/echo", echo)
	http.HandleFunc("/api/books", api.BookHandleFunc)
	http.HandleFunc("/api/books/", api.BookHandleFunc)
	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func index(rw http.ResponseWriter, rq *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "Hello World - Cloud Native Go")
	fmt.Fprintln(rw, "Version: 2")
}

func echo(rw http.ResponseWriter, rq *http.Request) {
	message := rq.URL.Query()["message"][0]
	rw.Header().Add("Content-Type", "text/plain")
	fmt.Fprintln(rw, message)
}
