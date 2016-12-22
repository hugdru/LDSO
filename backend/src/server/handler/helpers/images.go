package helpers

import (
	"net/http"
	"io/ioutil"
	"errors"
	"strconv"
)

const MaxImageFileSize = 16 << 20


var acceptedImageTypes = map[string]bool {
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
		return nil, "", errors.New("image is too big, max: "+strconv.Itoa(maxImageSize)+" bytes")
	}

	return imageBytes, imageMime, nil
}

