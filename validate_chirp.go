package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) validateChipr(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}
	type returnVals struct {
		// the key will be the name of struct field unless you give it an explicit JSON tag
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)

		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	respBody := returnVals{
		Body: filter(params.Body),
	}
	respondWithJSON(w, 200, respBody)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		// the key will be the name of struct field unless you give it an explicit JSON tag
		Error string `json:"error"`
	}
	dat, err := json.Marshal(returnError{
		Error: msg,
	})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
	return
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
}

func filter(msg string) string {
	notAllowedWords := map[string]int{
		"kerfuffle": 1,
		"sharbert":  2,
		"fornax":    3,
	}
	words := strings.Split(msg, " ")
	for i := 0; i < len(words); i++ {
		_, ok := notAllowedWords[strings.ToLower(words[i])]
		if ok {
			words[i] = "****"
		}
	}
	msg = strings.Join(words, " ")
	return msg
}
