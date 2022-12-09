package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"golang.org/x/net/webdav"
)

var Serve http.Handler

var DIRECTORY string

func main() {
	DIRECTORY := os.Getenv("DATAMODEL_DIRECTORY")
	DatamodelFileHandler := &webdav.Handler{
		FileSystem: webdav.Dir(DIRECTORY),
		LockSystem: webdav.NewMemLS(),
	}

	chiR := chi.NewRouter()
	chiR.Handle("/datamodel/{filename}", DatamodelFileHandler)
	chiR.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Reached Index Endpoint"))
	})

	log.Print("Server started on localhost, use /datamodel for uploading files / for reaching basic server")
	Serve = chiR
	log.Fatal(http.ListenAndServe(":80", Serve))
}
