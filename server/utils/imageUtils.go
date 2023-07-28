package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"os"
)

func StoreImageInNewFile(image []byte, path string, id int32, useCompression bool) string {
	var resPath = ""
	if len(image) > 0 {
		path = createPath(path, id)
		storeImageBytesAtPath(path, image, useCompression, id)
	}
	return resPath
}

func createPath(input string, id int32) string {
	path := fmt.Sprintf("%s%s", "./Storage/plants/", input)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.WithError(err).Errorf("Directory with path %s could not be created successfully", path)
			return ""
		}
	}
	return fmt.Sprintf("%s/%d", path, id)
}

// StoreImageBytesAtPath
// stores an image and returns teh path to this image.
// if needed, a new directory will be created and the path is extended until it is unique
func storeImageBytesAtPath(path string, i []byte, useCompression bool, id int32) {
	img, _, _ := image.Decode(bytes.NewReader(i))
	saveCompressedWithMaxSize(path, ".jpg", img, 1024*1024, len(i), useCompression)
	if id == 0 {
		saveCompressedWithMaxSize(path, "_thumbnail.jpg", img, 64*1024, len(i), useCompression)
	}
}

func saveCompressedWithMaxSize(path string, pathSuffix string, img image.Image, maxSize int32, imgSize int, useCompression bool) {
	out, errCreate := os.Create(fmt.Sprintf("%s%s", path, pathSuffix))
	if errCreate != nil {
		log.WithError(errCreate).Error("Error while creating a new file on the path: ", fmt.Sprintf("%s%s", path, pathSuffix))
		return
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing the file.")
		}
	}(out)

	if useCompression {
		scaleFactor := float64(maxSize) / float64(imgSize)
		newWidth := int(float64(img.Bounds().Dx()) * scaleFactor)
		newHeight := int(float64(img.Bounds().Dy()) * scaleFactor)
		img = imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	}

	errFile := jpeg.Encode(out, img, nil)
	if errFile != nil {
		log.WithError(errFile).Error("Error while encoding on the path: ", path)
		return
	}
}

func LoadImageBytesFromPath(path string) []byte {
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
