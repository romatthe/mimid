package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	"github.com/romatthe/mimid/resources"
	"log"
	"net/http"
)

func main() {

	// Let's fetch a DB instance
	db := NewDatabase()
	// Let's close it up when we exit the app
	defer db.Close()

	// Create a new httprouter
	router := httprouter.New()

	// Available routes
	router.GET("/", helloWorldHandler(db))
	router.GET("/hello/:name", helloYou)
	router.GET("/param", helloParams)

	// Songs
	router.GET("/songs", resources.GetSongs)

	// Stars HTTP Server and wire up the router
	http.ListenAndServe(":8080", router)
}

// NewDatabase instantiante a new Bolt.DB
func NewDatabase() *bolt.DB {
	// Open the mimid.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("mimid.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
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
func helloWorldHandler(db *bolt.DB) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		fmt.Fprintf(w, "Hi there my dear %s!", param.ByName("name"))
	})
}
