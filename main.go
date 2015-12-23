package main

import (
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	//"github.com/HouzuoGuo/tiedot/dberr"
	"github.com/julienschmidt/httprouter"
	"github.com/romatthe/mimid/resources"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	// Let's fetch a DB instance
	db := NewDatabase()
	// Let's close it up when we exit the app
	defer db.Close()

	// Create a new httprouter
	router := httprouter.New()

	// Available routes
	router.GET("/", homeHandler)
	router.GET("/hello/:name", helloYou)
	router.GET("/param", helloParams)

	// Songs
	router.GET("/songs", resources.GetSongs)

	// Upload experiment
	router.POST("/upload", uploadHandler)

	// Static file handler.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Stars HTTP Server and wire up the router
	http.ListenAndServe(":8080", router)
}

// NewDatabase instantiante a new tiedot db.DB
func NewDatabase() *db.DB {
	// Open the mimid.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := db.OpenDB("mimid.db")
	if err != nil {
		panic(err)
	}
	if err := db.Create("Feeds"); err != nil {
		panic(err)
	}
	return db
}

// Canonical Hello World function
func helloWorld(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!")
}

// Hello You!
func helloYou(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	fmt.Fprintf(w, "Hello %s!", param.ByName("name"))
}

// Hello Params!
func helloParams(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprint(w, r.Form.Get("floep"))
}

// Returns a handler for GET "/" and injects the DB instance
func helloWorldHandler(db *db.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		fmt.Fprintf(w, "Hi there my dear %s!", param.ByName("name"))
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	http.ServeFile(w, r, "index.html")
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
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
}
