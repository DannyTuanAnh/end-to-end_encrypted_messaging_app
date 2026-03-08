package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/h2non/filetype"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var (
	allowExts = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	allowMimeTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	allowFormats = map[string]bool{
		"jpeg": true,
		"png":  true,
		"jpg":  true,
	}
)

const maxFileSize = 5 << 20 // 5 MB

func ValidateAndSaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// 1. Validate filename - chặn path traversal và các ký tự đặc biệt
	// Chặn path separator (cả / và \)
	if strings.ContainsAny(fileHeader.Filename, "/\\") {
		return "", errors.New("filename contains path separator characters")
	}

	// Chặn path traversal
	if strings.Contains(fileHeader.Filename, "..") {
		return "", errors.New("filename contains path traversal sequences")
	}

	// Chặn các ký tự đặc biệt khác
	if strings.ContainsAny(fileHeader.Filename, `<>:"|?*`) {
		return "", errors.New("filename contains invalid special characters")
	}

	// 2. Check extension in file name
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowExts[ext] {
		return "", fmt.Errorf("unsupported file extension: (%s)", ext)
	}

	// 3. Check file size
	if fileHeader.Size > maxFileSize {
		return "", errors.New("File is too large (less than 5 MB)")
	}

	// 4. Open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf(`unable to open file: %v`, err)
	}
	defer file.Close()

	// 5. Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf(`unable to read file: %v`, err)
	}

	// 6. Validate image content
	kind, err := filetype.Match(fileBytes)
	if err != nil {
		return "", fmt.Errorf("cannot detect file type: %v", err)
	}

	// Check nếu không detect được định dạng
	if kind == filetype.Unknown {
		return "", errors.New("unknown file type")
	}

	// 7. Validate MIME type
	if !allowMimeTypes[kind.MIME.Value] {
		return "", fmt.Errorf("unsupported MIME type: %s (detected: %s)", kind.MIME.Value, kind.Extension)
	}

	// 8. Validate image true
	// This step decodes the image to ensure it is a valid image format.
	// If the image cannot be decoded, it is likely corrupted or not a valid image.
	// If the image is valid, it returns nil.
	// If the image is corrupted or invalid, it returns an error.
	// This step is important to ensure that the uploaded file is a valid image.
	if err := validateImageTrue(fileBytes); err != nil {
		return "", fmt.Errorf("corrupted or invalid image file: %v", err)
	}

	// Change file name
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create folder
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("cannot create upload folder")
	}

	//uploadDir "./uploads" + filename "abc.jpg"
	savePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := saveFile(fileHeader, savePath); err != nil {
		return "", fmt.Errorf("cannot save uploaded file: %v", err)
	}

	return filename, nil
}

// validateImageTrue func checks if the image data is a valid image format.
// It decodes the image and checks if the format is allowed.
// If the image cannot be decoded or the format is not allowed, it returns an error.
// If the image is valid, it returns nil.
func validateImageTrue(data []byte) error {
	_, format, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return fmt.Errorf("cannot decode image: %v", err)
	}

	if !allowFormats[format] {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	return nil
}

// saveFile func saves the uploaded file to the specified destination.
func saveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		os.Remove(destination)
		return err
	}

	return err
}
