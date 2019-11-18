package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const apiKey = "fb2cf5ec9500b1a5df159b1dbd5a29553b098a57283ba4ca22d28c788212d20e"
const baseURL = "https://arcentry.com/api/v1/"

type Document struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type DocumentById struct {
	DocId string `json:"docId"`
}

func getAllDocuments() ([]Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, baseURL+"doc", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	docs := make([]Document, 0)

	if err := json.Unmarshal(body, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func getIdFromName(name string) (string, error) {
	docs, err := getAllDocuments()
	if err != nil {
		return "", err
	}

	for _, doc := range docs {
		if doc.Title == name {
			return doc.Id, nil
		}
	}

	return "", errors.New("Document not found")
}

func createStaticEmbed(id string) ([]byte, error) {
	client := &http.Client{}

	doc := &DocumentById{DocId: id}

	json, err := json.Marshal(doc)

	req, err := http.NewRequest(http.MethodPost, baseURL+"embed/create-static", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func getEmbedByName(name string) ([]byte, error) {
	id, err := getIdFromName(name)
	if err != nil {
		return nil, err
	}

	return createStaticEmbed(id)
}
