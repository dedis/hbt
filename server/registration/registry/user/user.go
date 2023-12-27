package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.dedis.ch/hbt/server/registration/database"
	"go.dedis.ch/hbt/server/registration/registry"
)

var userDb database.Database

// RegisterDb registers the database for the user service
func RegisterDb(db database.Database) {
	userDb = db
}

// CreateDocument translates the http request to create a new document in the database
func CreateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal(err)
	}

	name := r.FormValue("name")
	passport := r.FormValue("passport")
	role, err := strconv.ParseUint(r.FormValue("role"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	picture, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println(fileHeader)

	picData := make([]byte, fileHeader.Size)
	picture.Read(picData)

	hash := r.FormValue("hash")

	regData := registry.RegistrationData{
		Name:     name,
		Passport: passport,
		Role:     uint(role),
		Picture:  picData,
		Hash:     []byte(hash),
	}

	docId, err := userDb.Create(regData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(docId)
	fmt.Println(docId)
}

// GetDocument translates the http request to get a document from the database
func GetDocument(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing id"))
	}

	hash := r.FormValue("hash")

	regId := registry.RegistrationId{
		Id: []byte(id),
	}

	data, err := userDb.Read(regId, []byte(hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(data)
}

// DeleteDocument translates the http request to delete a document in the database
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing id"))
	}

	hash := r.FormValue("hash")

	regId := registry.RegistrationId{
		Id: []byte(id),
	}

	err := userDb.Delete(regId, []byte(hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
