package api

import (
	"bufio"
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
	out, err := convert(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Here we need a valid json to send to arcentry, currently we only got an
	// already transformed json :-P
	res, err := getEmbedByName("consul")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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
	// translate document
	documentWithConsul, err := config.AddConsul(document)
	if err != nil {
		return
	}

	// getting the output hcl/json
	var b bytes.Buffer
	payload := bufio.NewWriter(&b)
	err = config.Write(payload, ext, documentWithConsul)
	if err != nil {
		err = fmt.Errorf("Error generating the consul output file: %w", err)
		return
	}

	// ------------ Image fetching with json only ----------------
	// currently we support bot json and hcl , but for arcentry integration it is only json based

	getJson := func(doc *config.Root) []byte {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		config.WriteJSON(writer, doc)
		return buf.Bytes()
	}

	// Would be bette just to check if it is hcl to do the conversion
	// and just calculate one, but for demo, let's do everything
	initial, final := getJson(document), getJson(documentWithConsul)
	// by now just adding
	// Polo : These are the jsons to plot with Arcentry
	out.Images = []*File{
		&File{
			Extension: ".json", // redudant
			Payload:   initial,
		},
		&File{
			Extension: ".json", // redundant
			Payload:   final,
		},
	}

	// -----------------------------------------------------------

	out.Code = http.StatusOK
	//redudant
	out.Consulfile = &File{
		Extension: input.Extension,
		Payload:   b.Bytes(),
	}

	return
}
