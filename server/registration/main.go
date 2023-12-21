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

func newApp() *application {
	userRouter := mux.NewRouter()
	userRouter.HandleFunc("/document", user.CreateDocument).Methods("POST")
	userRouter.HandleFunc("/document/{id}", user.GetDocument).Methods("GET")

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin/document/{id}", admin.GetDocument).Methods("GET")
	adminRouter.HandleFunc("/admin/document/{id}", admin.UpdateDocument).Methods("PUT")
	adminRouter.HandleFunc("/admin/document/{id}", admin.DeleteDocument).Methods("DELETE")

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
