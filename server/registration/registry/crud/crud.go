package crud

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go.dedis.ch/hbt/server/registration/database"
	"go.dedis.ch/hbt/server/registration/registry"
)

// CreateDocument translates the http request to create a new document in the database
func CreateDocument(w http.ResponseWriter, r *http.Request, db database.Database) {
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
	picture, fileHeader, err := r.FormFile("portrait")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	picData := make([]byte, fileHeader.Size)
	_, err = picture.Read(picData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	regData := &registry.RegistrationData{
		Name:       name,
		Passport:   passport,
		Role:       role,
		Picture:    picData,
		Registered: false,
	}

	registrationID, err := db.Create(regData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(registrationID)
	if err != nil {
		log.Println(err)
	}
	log.Println(registrationID)
}

// GetDocument translates the http request to get a document from the database
func GetDocument(w http.ResponseWriter, r *http.Request, db database.Database) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("missing id"))
		if err != nil {
			log.Println(err)
		}
	}

	hash := r.FormValue("hash")

	registrationID := registry.RegistrationID{
		ID: []byte(id),
	}

	data, err := db.Read(registrationID, []byte(hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(data)
	if err != nil {
		log.Println(err)
	}
	log.Println(data)
}

// UpdateDocument translates the http request to update a document in the database
func UpdateDocument(w http.ResponseWriter, r *http.Request, db database.Database, registered bool) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("missing id"))
		if err != nil {
			log.Println(err)
		}
	}

	registrationID := registry.RegistrationID{
		ID: []byte(id),
	}

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
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	log.Println(fileHeader)

	picData := make([]byte, fileHeader.Size)
	_, err = picture.Read(picData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	regData := &registry.RegistrationData{
		Name:     name,
		Passport: passport,
		Role:     role,
		Picture:  picData,
	}

	err = db.Update(registrationID, regData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(registrationID)
	if err != nil {
		log.Println(err)
	}
	log.Println(registrationID)
}

// DeleteDocument translates the http request to delete a document in the database
func DeleteDocument(w http.ResponseWriter, r *http.Request, db database.Database) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("missing id"))
		if err != nil {
			log.Println(err)
		}
	}

	hash := r.FormValue("hash")

	registrationID := registry.RegistrationID{
		ID: []byte(id),
	}

	err := db.Delete(registrationID, []byte(hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
