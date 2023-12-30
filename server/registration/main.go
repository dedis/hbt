package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.dedis.ch/hbt/server/registration/config"
	"go.dedis.ch/hbt/server/registration/database/mongodb"
	"go.dedis.ch/hbt/server/registration/registry/admin"
	"go.dedis.ch/hbt/server/registration/registry/user"
)

// curl -F "name='John Doe'" -F "passport=12XY456789" -F "role=0" -F "image=@./picture.png" -F "registered=false" localhost:80/document

// POST /document HTTP/1.1
// Host: localhost:3000
// User-Agent: curl/7.81.0
// Accept: */*
// Content-Length: 8722
// Content-Type: multipart/form-data; boundary=------------------------a7b1b827961ef70e
//
// --------------------------a7b1b827961ef70e
// Content-Disposition: form-data; name="name"
//
// 'John Doe'
// --------------------------a7b1b827961ef70e
// Content-Disposition: form-data; name="passport"
//
// 12XY456789
// --------------------------a7b1b827961ef70e
// Content-Disposition: form-data; name="role"
//
// 0
// --------------------------a7b1b827961ef70e
// Content-Disposition: form-data; name="image"; filename="picture.png"
// Content-Type: image/png
//
// �PNG

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
	userRouter.HandleFunc("/document", user.DeleteDocument).Methods("DELETE")

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin/document", admin.GetDocument).Methods("GET")
	adminRouter.HandleFunc("/admin/document", admin.UpdateDocument).Methods("PUT")
	adminRouter.HandleFunc("/admin/document", admin.DeleteDocument).Methods("DELETE")

	userDb := mongodb.NewUserDbAccess()
	user.RegisterDb(userDb)

	adminDb := mongodb.NewAdminDbAccess()
	admin.RegisterDb(adminDb)

	return &application{AdminRouter: adminRouter, UserRouter: userRouter}
}

func (a *application) start() {
	log.Println(fmt.Sprintf("Starting user server on port %s", config.AppConfig.UserServerPort))
	go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.UserServerPort),
		a.UserRouter))

	log.Println(fmt.Sprintf("Starting admin server on port %s", config.AppConfig.AdminServerPort))
	go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.AppConfig.AdminServerPort),
		a.AdminRouter))
}

func main() {
	config.LoadAppConfig()
	app := newApp()
	app.start()
}
