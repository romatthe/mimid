package main

import (
	"fmt"
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

// WorkerMusicUpload is a worker function called in order to process and organize new music files being uploaded
func WorkerMusicUpload(config Config, db *db.DB, uploads <-chan FileUpload, results chan<- FileUploadResult) {
	for upload := range uploads {
		filePath := path.Join(config.BaseMusicDir, upload.FileName)
		ioutil.WriteFile(filePath, upload.FileBuf, os.ModePerm)

		fmt.Printf("%+v\n", ParseMetaData(upload.UploadID, upload.FileBuf))

		results <- FileUploadResult{UploadID: upload.UploadID}
	}
}

// WorkerMusicUploadResult is a worker function for handling the results on
func WorkerMusicUploadResult(config Config, db *db.DB, uploadResults <-chan FileUploadResult) {
	for result := range uploadResults {
		fmt.Println(result.UploadID)
	}
}
