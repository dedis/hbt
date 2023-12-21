package admin

import (
	"net/http"

	"go.dedis.ch/hbt/server/registration/database"
)

var adminDb database.Database

func RegisterDb(db database.Database) {
	adminDb = db
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
}
