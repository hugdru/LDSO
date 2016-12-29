package helpers

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

const MaxImageFileSize = 16 << 20

var acceptedImageTypes = map[string]bool{
	"image/gif":  true,
	"image/png":  true,
	"image/jpeg": true,
}

func ReadImage(r *http.Request, imageName string, maxImageSize int) ([]byte, string, error) {
	mFile, mFileHeader, err := r.FormFile(imageName)
	if err == http.ErrMissingFile || err != nil {
		return nil, "", err
	}
	defer mFile.Close()

	imageContentType := GetContentType(mFileHeader.Header.Get("Content-type"))
	if !acceptedImageTypes[imageContentType] {
		return nil, "", errors.New("Gif, jpeg or png image")
	}
	mimetypeBuffer := make([]byte, 512)
	if _, err = mFile.Read(mimetypeBuffer); err != nil {
		return nil, "", errors.New("Gif, jpeg or png image")
	}
	if imageContentType != http.DetectContentType(mimetypeBuffer) {
		return nil, "", errors.New("Gif, jpeg or png image")
	}

	imageMime := imageContentType
	imageBytes, err := ioutil.ReadAll(mFile)
	if err != nil {
		return nil, "", err
	}

	if len(imageBytes) > maxImageSize {
		return nil, "", errors.New("image is too big, max: " + strconv.Itoa(maxImageSize) + " bytes")
	}

	return imageBytes, imageMime, err
}

func ReadImageBase64(image string, maxImageSize int) ([]byte, string, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return nil, "", err
	}

	if len(imageBytes) > maxImageSize {
		return nil, "", errors.New("image is too big, max: " + strconv.Itoa(maxImageSize) + " bytes")
	}

	mimetypeBuffer := make([]byte, 512)
	if copy(mimetypeBuffer, imageBytes) != 512 {
		return nil, "", errors.New("Expecting at least 512 bytes")
	}

	imageMime := http.DetectContentType(mimetypeBuffer)
	if !acceptedImageTypes[imageMime] {
		return nil, "", errors.New("Gif, jpeg or png image")
	}

	return imageBytes, imageMime, err
}
