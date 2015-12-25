package main

import (
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/julienschmidt/httprouter"
	. "github.com/romatthe/mimid/mimid"
	"net/http"
	"os"
	"os/user"
	"path"
)

func main() {

	// Bootstrap the application and configure the dependencies
	config, db := Startup()

	// Set up the channels on which to perform work
	fileUploads := make(chan FileUpload, 100)
	fileUploadResults := make(chan FileUploadResult, 100)

	// Create a bunch of nice workers for the fileUploads channel
	// TODO: Fetch the amount of workers from Config?
	for w := 0; w <= 50; w++ {
		go WorkerMusicUpload(config, db, fileUploads, fileUploadResults)
	}

	go WorkerMusicUploadResult(config, db, fileUploadResults)

	// Create a new httprouter
	router := httprouter.New()

	// Home and Static files
	router.GET("/", HandlerHome(config, db))
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	// Upload experiment
	router.POST("/upload", HandlerUpload(config, db, fileUploads))

	// Stars HTTP Server and wire up the router
	http.ListenAndServe(":8080", router)
}

// Startup bootstraps the application
func Startup() (Config, *db.DB) {
	// Fetch the current user context
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Create a standard Config object
	// TODO: Read from file, maybe HCL?
	baseDir := path.Join(usr.HomeDir, ".mimid")
	config := Config{BaseConfigDir: baseDir, BaseMusicDir: path.Join(baseDir, "music"), BaseDBDir: path.Join(baseDir, "mimid.db")}

	// TODO: Dev stuff, cleaning all files, remove this later
	_ = os.RemoveAll(config.BaseMusicDir)
	// Base config dir & Base music dir
	if err := os.MkdirAll(config.BaseConfigDir, 0766); err != nil {
		panic(err.Error())
	}
	if err := os.MkdirAll(config.BaseMusicDir, 0766); err != nil {
		panic(err.Error())
	}

	// Open/Create a tiedot DB
	db, err := db.OpenDB(config.BaseDBDir)
	if err != nil {
		panic(err)
	}

	return config, db
}
