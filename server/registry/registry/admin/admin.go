package admin

import (
	"net/http"

	"go.dedis.ch/hbt/server/registry/database"
	"go.dedis.ch/hbt/server/registry/registry/crud"
)

var adminDB database.Database

// RegisterDB registers the database for the admin service
func RegisterDB(db database.Database) {
	adminDB = db
}

// GetDocument translates the http request to get a document from the database
func GetDocument(w http.ResponseWriter, r *http.Request) {
	crud.GetDocument(w, r, adminDB)
}

// UpdateDocument translates the http request to update a document in the database
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	crud.UpdateDocument(w, r, adminDB, true)
}

// DeleteDocument translates the http request to delete a document from the database
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	crud.DeleteDocument(w, r, adminDB)
}
