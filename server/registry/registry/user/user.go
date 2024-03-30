package user

import (
	"net/http"

	"go.dedis.ch/hbt/server/registry/database"
	"go.dedis.ch/hbt/server/registry/registry/crud"
)

var userDB database.Database

// RegisterDB registers the database for the user service
func RegisterDB(db database.Database) {
	userDB = db
}

// CreateDocument translates the http request to create a new document in the database
func CreateDocument(w http.ResponseWriter, r *http.Request) {
	crud.CreateDocument(w, r, userDB)
}

// GetDocument translates the http request to get a document from the database
func GetDocument(w http.ResponseWriter, r *http.Request) {
	crud.GetDocument(w, r, userDB)
}

// UpdateDocument translates the http request to update a document in the database
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	crud.UpdateDocument(w, r, userDB, false)
}
