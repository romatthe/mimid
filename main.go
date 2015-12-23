package main

import (
	"github.com/HouzuoGuo/tiedot/db"
	//"github.com/HouzuoGuo/tiedot/dberr"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"os/user"
	"path"
)

func main() {

	// Bootstrap the application and configure the dependencies
	config, db := Startup()

	// Create a new httprouter
	router := httprouter.New()

	// Available routes
	router.GET("/", handlerHome(config, db))

	// Upload experiment
	router.POST("/upload", handlerUpload(config, db))

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
