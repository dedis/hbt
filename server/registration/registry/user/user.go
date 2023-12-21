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

	fmt.Println(fileHeader)

	picData := make([]byte, fileHeader.Size)
	picture.Read(picData)

	regData := registry.RegistrationData{
		Name:     name,
		Passport: passport,
		Role:     uint(role),
		Picture:  picData,
	}

	docId, err := userDb.Create(regData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	response := registry.RegistrationId{Id: docId}

	json.NewEncoder(w).Encode(response)
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
}
