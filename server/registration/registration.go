package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.dedis.ch/hbt/server/registration/config"
	"go.dedis.ch/hbt/server/registration/database/mongodb"
	"go.dedis.ch/hbt/server/registration/registry/admin"
	"go.dedis.ch/hbt/server/registration/registry/user"
)

// curl -F "name='John Doe'" -F "passport=12XY456789" -F "role=0" -F "image=@test/passport.jpg" -F "registered=false" localhost:3000/document

// application defines the application instance
type application struct {
	UserRouter  *mux.Router
	AdminRouter *mux.Router
}

// newApp creates a new application instance
// it creates the user and admin router
// and registers the handlers
// the routes used are pretty simple:
// - /document for the user
// - /admin/document for the admin
// to create a new document, the user sends a POST request to /document
// with the following form data:
// - name: string
// - passport: string
// - role: uint
// - image: file
// the response is a json string containing the document id
// to get a document, the user sends a GET request to /document
// with the following query parameter:
// - id: string
// to delete a document, the user sends a DELETE request to /document
// with the following query parameter:
// - id: string
// to get a document, the admin sends a GET request to /admin/document
// with the following query parameter:
// - id: string
// to update a document, the admin sends a PUT request to /admin/document
// with the following query parameter:
// - id: string
// and the following form data:
// - name: string
// - passport: string
// - role: uint
// - registered: bool
func newApp() *application {
	userRouter := mux.NewRouter()
	userRouter.HandleFunc("/document", user.CreateDocument).Methods("POST")
	userRouter.HandleFunc("/document", user.GetDocument).Methods("GET")
	userRouter.HandleFunc("/document", user.UpdateDocument).Methods("PUT")

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin/document", admin.GetDocument).Methods("GET")
	adminRouter.HandleFunc("/admin/document", admin.UpdateDocument).Methods("PUT")
	adminRouter.HandleFunc("/admin/document", admin.DeleteDocument).Methods("DELETE")

	db, err := mongodb.NewDBAccess()
	if err != nil {
		log.Fatal(err)
	}
	user.RegisterDB(db)

	db, err = mongodb.NewDBAccess()
	if err != nil {
		log.Fatal(err)
	}
	admin.RegisterDB(db)

	return &application{AdminRouter: adminRouter, UserRouter: userRouter}
}

func (a *application) start() {
	log.Printf("Starting user server on port %s", config.AppConfig.UserServerPort)
	s := &http.Server{
		Addr:              fmt.Sprintf(":%v", config.AppConfig.UserServerPort),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           a.UserRouter,
	}
	go log.Fatal(s.ListenAndServe())

	log.Printf("Starting admin server on port %s", config.AppConfig.AdminServerPort)
	s = &http.Server{
		Addr:              fmt.Sprintf(":%v", config.AppConfig.AdminServerPort),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           a.AdminRouter,
	}
	go log.Fatal(s.ListenAndServe())
}

func main() {
	config.LoadAppConfig()
	app := newApp()
	app.start()
}
