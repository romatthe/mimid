package main

import (
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	//"github.com/HouzuoGuo/tiedot/dberr"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func handlerHome(config Config, db *db.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		http.ServeFile(w, r, "index.html")
	})
}

// This is where the upload action happens.
func handlerUpload(config Config, db *db.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		if err := r.ParseMultipartForm(1 * 1024 * 1024); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		for key, value := range r.MultipartForm.Value {
			fmt.Fprintf(w, "%s:%s ", key, value)
			log.Printf("%s:%s", key, value)
		}

		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				file, _ := fileHeader.Open()
				path := fmt.Sprintf("files/%s", fileHeader.Filename)
				buf, _ := ioutil.ReadAll(file)
				fmt.Fprintf(w, "files/%s", fileHeader.Filename)
				ioutil.WriteFile(path, buf, os.ModePerm)
			}
		}
	})
}
