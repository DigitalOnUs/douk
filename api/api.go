package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/DigitalOnUs/inotx/config"
)

//////// models //////////

var (
	ErrEmptyInputFile      = errors.New("Empty file to consulize")
	ErrorUnsupportedFormat = errors.New("Not supported extension")
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

	// basic validations
	_, err := convert(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", input)
}

func convert(input File) (out Response, err error) {
	out = Response{}

	if input.Extension != ".hcl" && input.Extension != ".json" {
		err = fmt.Errorf("%w : %s", ErrorUnsupportedFormat, input.Extension)
		return
	}

	ext := input.Extension[1:]
	defaultName := "inputDocument" + input.Extension

	reader := bytes.NewReader(input.Payload)
	document, err := config.Parse(reader, defaultName, ext)
	if err != nil {
		return
	}

	fmt.Printf("%+v\n", document)

	return
}
