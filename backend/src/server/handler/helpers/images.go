package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
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

func ReadImage(r *http.Request, imageName string, maxImageSize int) ([]byte, string, string, error) {
	mFile, mFileHeader, err := r.FormFile(imageName)
	if err == http.ErrMissingFile || err != nil {
		return nil, "", "", err
	}
	defer mFile.Close()

	imageContentType := GetContentType(mFileHeader.Header.Get("Content-type"))
	if !acceptedImageTypes[imageContentType] {
		return nil, "", "", errors.New("Gif, jpeg or png image")
	}
	mimetypeBuffer := make([]byte, 512)
	readBytes, err := mFile.Read(mimetypeBuffer)
	if err != nil {
		return nil, "", "", errors.New("Gif, jpeg or png image")
	}
	if readBytes != 512 {
		return nil, "", "", errors.New("Expecting at least 512 bytes")
	}

	if imageContentType != http.DetectContentType(mimetypeBuffer) {
		return nil, "", "", errors.New("Gif, jpeg or png image")
	}

	imageMime := imageContentType

	_, err = mFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, "", "", err
	}
	imageBytes, err := ioutil.ReadAll(mFile)
	if err != nil {
		return nil, "", "", err
	}

	if len(imageBytes) > maxImageSize {
		return nil, "", "", errors.New("image is too big, max: " + strconv.Itoa(maxImageSize) + " bytes")
	}

	imageHash := CalculateImageHash(imageBytes)

	return imageBytes, imageMime, imageHash, err
}

func ReadImageBase64(image string, maxImageSize int) ([]byte, string, string, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return nil, "", "", err
	}

	if len(imageBytes) > maxImageSize {
		return nil, "", "", errors.New("image is too big, max: " + strconv.Itoa(maxImageSize) + " bytes")
	}

	mimetypeBuffer := make([]byte, 512)
	if copy(mimetypeBuffer, imageBytes) != 512 {
		return nil, "", "", errors.New("Expecting at least 512 bytes")
	}

	imageMime := http.DetectContentType(mimetypeBuffer)
	if !acceptedImageTypes[imageMime] {
		return nil, "", "", errors.New("Gif, jpeg or png image")
	}

	imageHash := CalculateImageHash(imageBytes)

	return imageBytes, imageMime, imageHash, err
}

// sha256 to hex has 64 chars
const encodedImageHashSize = 64

func CalculateImageHash(imageBytes []byte) string {
	binaryHash := sha256.Sum256(imageBytes)
	encodedHash := hex.EncodeToString(binaryHash[:])
	return encodedHash
}

func ImageHashSizeIsValid(size int) bool {
	return size == encodedImageHashSize
}

func ImageCashingControlWriter(w http.ResponseWriter, r *http.Request, urlHash, imageHash, imageMimetype string, imageBytes []byte, redirectUrl string) error {
	if urlHash != imageHash {
		http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
		return nil
	}
	entityEtag := imageHash
	if getEtag := r.Header.Get("If-None-Match"); getEtag != "" {
		if getEtag == entityEtag {
			w.Header().Set("Cache-Control", "public, max-age=2592000")
			w.WriteHeader(http.StatusNotModified)
			return nil
		}
	}
	w.Header().Set("Content-type", imageMimetype)
	w.Header().Set("Content-Length", strconv.Itoa(len(imageBytes)))
	w.Header().Set("Cache-Control", "public, max-age=2592000")
	w.Header().Set("Etag", entityEtag)
	w.Write(imageBytes)
	return nil
}
