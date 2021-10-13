package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/corona10/goimagehash"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type Labels struct {
	Name       string
	Confidence float64
}

type PhotoMetaData struct {
	Name                string
	ParsedName          string
	PrefixName          string
	Artist              string
	CaptureTime         time.Time
	CaptureYear         string
	CaptureYearMonth    string
	CaptureYearMonthDay string
	UploadTime          time.Time
	Description         string
	Caption             string
	Height              int
	Width               int
	Sha256              string
	PerceptualHash      string
	// Classification      struct {
	// 	Labels []Labels
	// }
}

func getCleanExifValue(md *tiff.Tag) string {
	if md == nil {
		return ""
	}
	s := fmt.Sprintf("%v", md)

	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func populatePMD(filepath string) *PhotoMetaData {
	fileBytes, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	defer fileBytes.Close()

	x, err := exif.Decode(fileBytes)

	if err != nil {
		fmt.Print("should not get an error")
	}

	var pmd *PhotoMetaData = new(PhotoMetaData)

	pmd.UploadTime = time.Now().UTC()

	exifValueArtist, err := x.Get(exif.Artist)

	if err != nil {
		fmt.Print("error decoding the Artist")
	}

	pmd.Artist = getCleanExifValue(exifValueArtist)

	pmd.CaptureTime, err = x.DateTime()

	if err != nil {
		fmt.Print("error decodeing the time")
	}

	// `Format` and `Parse` use example-based layouts. Usually
	// you'll use a constant from `time` for these layouts, but
	// you can also supply custom layouts. Layouts must use the
	// reference time `Mon Jan 2 15:04:05 MST 2006` to show the
	// pattern with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006,
	// 15 for the hour, Monday for the day of the week, etc.
	pmd.CaptureYear = pmd.CaptureTime.Format("2006")
	pmd.CaptureYearMonth = pmd.CaptureTime.Format("2006-01")
	pmd.CaptureYearMonthDay = pmd.CaptureTime.Format("2006-01-02")

	exifValueDescription, _ := x.Get(exif.ImageDescription)

	pmd.Description = getCleanExifValue(exifValueDescription)

	// Need to rewind to the start of the file so I can get the dimensions (since they aren't in the EXIF)
	fileBytes.Seek(0, io.SeekStart)
	im, _, err := image.DecodeConfig(fileBytes)

	if err != nil {
		log.Printf("cannot open %s to get dimensions\n", filepath)
	} else {
		pmd.Width = im.Width
		pmd.Height = im.Height
	}

	// Need to rewind the file again to get the sha256
	fileBytes.Seek(0, io.SeekStart)
	h := sha256.New()
	if _, err := io.Copy(h, fileBytes); err != nil {
		panic(err)
	}
	pmd.Sha256 = fmt.Sprintf("%x", h.Sum(nil))

	// Perceptual hash (and yet another rewind of the file)
	fileBytes.Seek(0, io.SeekStart)
	img1, _ := jpeg.Decode(fileBytes)

	phash, _ := goimagehash.PerceptionHash(img1)
	pmd.PerceptualHash = fmt.Sprintf("%x", phash.GetHash())

	// Rekognition
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
	}

	svc := rekognition.New(sess)

	fileBytes.Seek(0, io.SeekStart)
	reader := bufio.NewReader(fileBytes)
	content, _ := ioutil.ReadAll(reader)

	inputRkg := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: content,
		},
	}

	result, err := svc.DetectLabels(inputRkg)

	if err != nil {
		log.Printf("error with DetectLabels %v\n", err)
	}

	for _, lab := range result.Labels {
		l := Labels{*lab.Name, *lab.Confidence}
		pmd.Classification.Labels = append(pmd.Classification.Labels, l)
	}

	return pmd
}
