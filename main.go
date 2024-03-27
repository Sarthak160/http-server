package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var books []Book

func main() {
    http.HandleFunc("/books", handleBooks)
    log.Fatal(http.ListenAndServe(":8085", nil))
}

func handleBooks(w http.ResponseWriter, r *http.Request) {
    fmt.Println("INSIDE HANDLE BOOKS !!!!!")
	switch r.Method {
    case http.MethodGet:
        getBooks(w, r)
    case http.MethodPost:
        createBook(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
    }
}

func getBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
    var newBook Book
    err := json.NewDecoder(r.Body).Decode(&newBook)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid request payload"))
        return
    }

    books = append(books, newBook)
    w.WriteHeader(http.StatusCreated)
}
