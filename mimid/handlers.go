package mimid

import (
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
)

// HandlerHome returns an httprouter.Handle to serve the main index.html
func HandlerHome(config Config, db *db.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		http.ServeFile(w, r, "index.html")
	})
}

// HandlerUpload returns an httrouter.Handle to manage the MultipartForm file upload
func HandlerUpload(config Config, db *db.DB, uploads chan<- FileUpload) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		if err := r.ParseMultipartForm(1 * 1024 * 1024); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		// We want to enumerate all files and put them on a channel for workers to be picked up
		for _, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				// Let's open the file header so wel can pull out the name and byte buf
				file, _ := fileHeader.Open()
				buf, _ := ioutil.ReadAll(file)

				// Let's put the file on a worker channel so it can be properly processed
				upload := FileUpload{
					UploadID: HashFileName(fileHeader.Filename),
					FileBuf:  buf,
					FileName: fileHeader.Filename,
				}
				uploads <- upload
			}
		}

		w.WriteHeader(http.StatusAccepted)
	})
}
