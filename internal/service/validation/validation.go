package validation

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^.+@.+$`)
	return re.MatchString(email)
}

func ValidateLogin(login string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-_.]{3,20}[A-Za-z0-9]$`)
	return re.MatchString(login)
}

func ValidatePassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=-]{8,16}$`)
	return re.MatchString(password)
}

func ValidateName(name string) bool {
	re := regexp.MustCompile(`^.{5,50}$`)
	return re.MatchString(name)
}

func ValidateImages(files []*multipart.FileHeader, maxSize int64, allowedMimeTypes []string, maxWidth, maxHeight int) error {
	for _, file := range files {
		if err := ValidateImage(file, maxSize, allowedMimeTypes, maxWidth, maxHeight); err != nil {
			return fmt.Errorf("file %s is invalid: %w", file.Filename, err)
		}
	}
	return nil
}

func ValidateImage(file *multipart.FileHeader, maxSize int64, allowedMimeTypes []string, maxWidth, maxHeight int) error {
	if file.Size > maxSize {
		return fmt.Errorf("file exceeds maximum size of %d bytes", maxSize)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return fmt.Errorf("could not read file buffer: %v", err)
	}

	mimeType := http.DetectContentType(buffer)
	allowed := false
	for _, t := range allowedMimeTypes {
		if mimeType == t {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type %s is not allowed", mimeType)
	}

	src.Seek(0, 0)
	var img image.Image
	switch {
	case strings.HasSuffix(mimeType, "jpeg"):
		img, err = jpeg.Decode(src)
	case strings.HasSuffix(mimeType, "png"):
		img, err = png.Decode(src)
	case strings.HasSuffix(mimeType, "jpg"):
		img, err = jpeg.Decode(src)
	default:
		return errors.New("unsupported image format")
	}
	if err != nil {
		return fmt.Errorf("could not decode image: %v", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width > maxWidth || height > maxHeight {
		return fmt.Errorf("image resolution exceeds maximum allowed size of %dx%d", maxWidth, maxHeight)
	}

	return nil
}
