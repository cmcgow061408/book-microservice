package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookToJSON(t *testing.T) {
	book := Book{Title: "Cloud Native Go", Author: "M.-L. Reimer", ISBN: "0123455677"}
	rs := book.toJSON()
	assert.Equal(t, "{\"title\":\"Cloud Native Go\",\"author\":\"M.-L. Reimer\",\"isbn\":\"0123455677\"}", string(rs), "Book - JSON Marshalling is wrong.")
}

func TestBookFromJSON(t *testing.T) {
	data := []byte("{\"title\":\"Cloud Native Go\",\"author\":\"M.-L. Reimer\",\"isbn\":\"0123455677\"}")
	book := Book{Title: "Cloud Native Go", Author: "M.-L. Reimer", ISBN: "0123455677"}

	assert.Equal(t, book, fromJSON(data), "Book - JSON Unmarshalling is wrong.")

}
