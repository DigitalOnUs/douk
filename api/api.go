package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//////// models //////////

var (
	ErrEmptyInputFile = errors.New("Empty file to consulize")
)

//File struct
type File struct {

	// extension of the file
	Extension string `json:"extension,omitempty"`

	// content of the file
	Payload []byte `json:"payload,omitempty"`
}

//Response default
type Response struct {
	Consulfile *File   `json:"consulfile,omitempty"`
	Images     []*File `json:"images,omitempty"`
	Code       int32   `json:"code,omitempty"`
	Message    string  `json:"message,omitempty"`
}

///// end models ///////

//Consulize add the values
func Consulize(w http.ResponseWriter, r *http.Request) {
	var input File

	decoder := json.NewDecoder(r.Body)
	// validations
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(input.Payload) < 1 {
		http.Error(w, ErrEmptyInputFile.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", input)
}
