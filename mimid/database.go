package mimid

import (
	"github.com/HouzuoGuo/tiedot/db"
)

// Database defines all database operations
type Database interface {
	GetSongs() []Track
}

// ApplicationDatabase exposes the application database
type ApplicationDatabase struct {
	DB *db.DB
}

// NewApplicationDatabase constructs a new ApplicationDatabase object
func NewApplicationDatabase(config Config) *ApplicationDatabase {
	// Open/Create a tiedot DB
	db, err := db.OpenDB(config.BaseDBDir)
	if err != nil {
		panic(err)
	}
	// Init the Tracks collection
	if err := db.Create("Tracks"); err != nil {
		panic(err)
	}

	return &ApplicationDatabase{DB: db}
}
