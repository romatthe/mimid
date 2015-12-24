package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/HouzuoGuo/tiedot/db"
	"io/ioutil"
	"os"
	"path"
)

// FileUpload represents an uploaded file that still needs to be processed by a worker
type FileUpload struct {
	UploadID string
	FileBuf  []byte
	FileName string
}

// FileUploadResult presents the result for an uploaded file with a specifc id
type FileUploadResult struct {
	UploadID string
}

// HashFileName takes a string containing a filename and hashes it to uniquely identify it
// This is the name + extension only, no path
func HashFileName(fileName string) string {
	data := md5.Sum([]byte(fileName))
	return hex.EncodeToString(data[:])
}

// WorkerMusicUpload is a worker function called in order to process and organize new music files being uploaded
func WorkerMusicUpload(config Config, db *db.DB, uploads <-chan FileUpload, results chan<- FileUploadResult) {
	for upload := range uploads {
		filePath := path.Join(config.BaseMusicDir, upload.FileName)
		ioutil.WriteFile(filePath, upload.FileBuf, os.ModePerm)
		results <- FileUploadResult{UploadID: upload.UploadID}
	}
}
