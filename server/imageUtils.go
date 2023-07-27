package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"os"
)

func imageWrapper(image []byte, path string, id int32) string {
	var resPath = ""
	if len(image) > 0 {
		var resError error
		path := fmt.Sprintf("%s%s%s%d%s", "./Storage/rating/", path, "/", id, "/")
		resPath, resError = StoreImageBytesAtPath(path, image)

		if resError != nil {
			log.WithError(resError).Error("Error occurred while storing the image.")
		}
	}
	return resPath
}

// StoreImageBytesAtPath
// stores an image and returns teh path to this image.
// if needed, a new directory will be created and the path is extended until it is unique
func StoreImageBytesAtPath(path string, i []byte) (string, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.WithError(err).Errorf("Directory with path %s could not be created successfully", path)
			return "", nil
		}
	}

	img, _, _ := image.Decode(bytes.NewReader(i))
	resizedImage := imaging.Resize(img, 1280, 0, imaging.Lanczos)

	var opts jpeg.Options
	maxImageSize := 524288 // 0.55MB
	if len(i) > maxImageSize {
		opts.Quality = (maxImageSize / len(i)) * 100
	} else {
		opts.Quality = 100 // if image small enough use it directly
	}

	var imgPath = fmt.Sprintf("%s%x.jpeg", path, md5.Sum(i))

	out, errFile := os.Create(imgPath)
	if errFile != nil {
		log.WithError(errFile).Error("Error while creating a new file on the path: ", path)
		return imgPath, errFile
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing the file.")
		}
	}(out)

	errFile = jpeg.Encode(out, resizedImage, &opts)
	return imgPath, errFile
}

func loadImageBytesFromPath(path string) []byte {
	if len(path) == 0 {
		return make([]byte, 0)
	}
	file, err := os.Open(path)

	if err != nil {
		log.WithError(err).Error("Error while opening image file with path: ", path)
		return nil
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("Unable to close the file for storing the image.")
		}
	}(file)

	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	imageAsBytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(imageAsBytes)
	if err != nil {
		log.WithError(err).Error("Error while trying to read image as bytes")
		return nil
	}
	return imageAsBytes
}
