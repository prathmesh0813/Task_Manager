package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// Validate avatar size and file extension
func ValidateAvatar(fileheader *multipart.FileHeader) (string, error) {

	if (fileheader.Size) > 100*1024 {

		return "", errors.New("File size is too large, max 100 kb is allowed")
	}

	fileExtension := strings.ToLower(filepath.Ext(fileheader.Filename))
	if fileExtension != ".jpg" && fileExtension != ".jpeg" {

		return fileExtension, errors.New("Invalid file type")
	}
	return fileExtension, nil
}
