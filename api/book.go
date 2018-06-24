package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//Book struct
type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description,omitempty"`
	//define the book
}

func (b Book) toJSON() []byte {
	toJSON, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return toJSON
}

func fromJSON(data []byte) Book {
	b := Book{}
	err := json.Unmarshal(data, &b)
	if err != nil {
		panic(err)
	}
	return b
}

//Books - A map containing books
var Books = map[string]Book{
	"0123455677": Book{Title: "Cloud Native Go", Author: "M.-L. Reimer", ISBN: "0123455677"},
	"0123455681": Book{Title: "The Hitchhiker's Guide to the Galaxy", Author: "Douglas Adams", ISBN: "0123455681"}}

//BookHandleFunc - Handles REST API calls for Books
func BookHandleFunc(rw http.ResponseWriter, rq *http.Request) {
	index := strings.LastIndex(rq.URL.Path, "/api/books/")
	isbn := ""
	if index > -1 {
		isbn = rq.URL.Path[len("/api/books/"):]
	}

	switch method := rq.Method; method {
	case http.MethodGet:
		if index > -1 {
			bk, found := GetBook(isbn)
			if found {
				rw.Header().Add("Content-Type", "application/json;charset=utf-8")
				rw.WriteHeader(http.StatusFound)
				rw.Write(bk.toJSON())
			} else {
				rw.WriteHeader(http.StatusNotFound)
			}
		} else {
			bks := AllBooks()
			byteBk, err := json.Marshal(bks)
			if err == nil {
				rw.Header().Add("Content-Type", "application/json;charset=utf-8")
				rw.WriteHeader(http.StatusOK)
				rw.Write(byteBk)
			} else {
				rw.WriteHeader(http.StatusSeeOther)
				rw.Write([]byte("Unable to return requested Books [JSON]."))
			}
		}
	case http.MethodPost:
		body, err := ioutil.ReadAll(rq.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		bk := fromJSON(body)
		isbn, created := CreateBook(bk)

		if len(isbn) == 0 {
			fmt.Printf("Book[%v] already existed within DB. Time: %v.", bk.ISBN, created)
			rw.WriteHeader(http.StatusConflict)
		} else {
			fmt.Printf("Added Book[%v] at %v.", isbn, created)
			rw.Header().Add("Location", "/api/books/"+isbn)
			rw.WriteHeader(http.StatusCreated)

		}
	case http.MethodPut:
		if index > -1 {
			body, err := ioutil.ReadAll(rq.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
			}
			bk := fromJSON(body)
			updated := UpdateBook(isbn, bk)
			if updated {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.WriteHeader(http.StatusNotFound)
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		if index > -1 {
			deleted := DeleteBook(isbn)
			if deleted {
				rw.WriteHeader(http.StatusOK)
			} else {
				rw.WriteHeader(http.StatusNotFound)
			}
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}

	default:
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Unsupported Resquest"))
	}
}

//AllBooks - Returns all books as slice from MAP of Books
func AllBooks() []Book {
	vs := make([]Book, 0, len(Books))

	for _, value := range Books {
		vs = append(vs, value)
	}
	return vs
}

//CreateBook - Adds book to Map of Books -> if an actual book and doesn't exist already within Map.
func CreateBook(b Book) (isbn string, created time.Time) {
	_, ok := GetBook(b.ISBN)
	if ok == false {
		//Add book because it doesn't already exists
		fmt.Printf("Adding Book: %v\n", b.ISBN)
		Books[b.ISBN] = b
		return b.ISBN, time.Now()
	}
	tm := time.Now()

	return "", tm
}

//GetBook - Returns a Book and a boolean representing if the book was found in Datasource
func GetBook(isbn string) (bk Book, found bool) {
	fmt.Printf("Finding Book associated with ISBN:%v", isbn)
	book, ok := Books[isbn]
	return book, ok
}

//UpdateBook - Updates a Books values based on ISBN identification.
func UpdateBook(isbn string, bk Book) (exists bool) {
	_, doesExist := GetBook(isbn)
	if doesExist {
		Books[isbn] = bk
	}

	return doesExist
}

//DeleteBook - Removes Book associated with ISBN identification.
func DeleteBook(isbn string) (deleted bool) {
	_, doesExist := GetBook(isbn)
	if doesExist {
		delete(Books, isbn)
		return true
	}

	return false
}
