package amd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/dhowden/tag"
)

type AudioMetaData struct {
	Name          string
	Artist        string    `json:"Artist"`
	Title         string    `json:"Title"`
	Genre         string    `json:"Genre"`
	Composer      string    `json:"Composer"`
	RunThrough    bool      `json:"RunThrough"`
	Piano         string    `json:"Piano"`
	Comments      string    `json:"Comments"`
	RecordingDate time.Time `json:"RecordingDate"`
	UploadDate    time.Time
	Sha256        string
}

func PopulateAMD(filepath string) *AudioMetaData {
	fileBytes, err := os.Open(filepath)
	if err != nil {
		log.Printf("err: %x", err)
	}

	defer fileBytes.Close()

	var amd *AudioMetaData = new(AudioMetaData)

	amd.UploadDate = time.Now().UTC()

	fileBytes.Seek(0, io.SeekStart)
	m, err := tag.ReadFrom(fileBytes)
	if err != nil {
		log.Printf("err: %x", err)
	}

	err = json.Unmarshal([]byte(m.Comment()), amd)
	if err != nil {
		log.Printf("err: %x", err)
	}

	// Need to rewind the file again to get the sha256
	fileBytes.Seek(0, io.SeekStart)
	h := sha256.New()
	if _, err := io.Copy(h, fileBytes); err != nil {
		log.Printf("err: %x", err)
	}
	amd.Sha256 = fmt.Sprintf("%x", h.Sum(nil))

	return amd
}
