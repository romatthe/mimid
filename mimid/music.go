package mimid

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/dhowden/tag"
)

// Track represents both the music and mimid metadata of a single track
type Track struct {
	TrackID     string
	Format      tag.Format
	FileType    tag.FileType
	Title       string
	Album       string
	Artist      string
	AlbumArtist string
	Composer    string
	Genre       string
	Year        int
	Track       int
	Disc        int
}

// HashFileName takes a string containing a filename and hashes it to uniquely identify it
// This is the name + extension only, no path
func HashFileName(fileName string) string {
	data := md5.Sum([]byte(fileName))
	return hex.EncodeToString(data[:])
}

// HashFileContent generates an MD5 hash based on the byte buffer of the file
func HashFileContent(fileContent []byte) string {
	data := md5.Sum(fileContent)
	return hex.EncodeToString(data[:])
}

// ParseMetaData tries to fill in the initial metadata as well as possible based on raw file data
func ParseMetaData(fileID string, fileContent []byte) Track {
	// Test to read some Metadata
	reader := bytes.NewReader(fileContent)
	meta, err := tag.ReadFrom(reader)
	if err != nil {
		panic(err)
	}

	// Get singular value from multi-value fields
	track, _ := meta.Track()
	disc, _ := meta.Disc()

	return Track{
		TrackID:     fileID,
		Format:      meta.Format(),
		FileType:    meta.FileType(),
		Title:       meta.Title(),
		Album:       meta.Album(),
		Artist:      meta.Artist(),
		AlbumArtist: meta.AlbumArtist(),
		Composer:    meta.Composer(),
		Genre:       meta.Genre(),
		Year:        meta.Year(),
		Track:       track,
		Disc:        disc,
	}
}
