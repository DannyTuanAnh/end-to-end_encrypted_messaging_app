package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"
)

// Helper: Tạo ảnh JPEG giả
func createFakeJPEG(t *testing.T) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// Helper: Tạo ảnh PNG giả
func createFakePNG(t *testing.T) []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// Helper: Tạo multipart.FileHeader từ bytes
func createFileHeader(filename string, data []byte) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("file", filename)
	part.Write(data)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(10 << 20)

	return form.File["file"][0]
}

// ✅ Helper MỚI (force filename, không normalize)
func createFileHeaderWithExactFilename(filename string, data []byte) *multipart.FileHeader {
	// Tạo FileHeader thủ công, không qua multipart parsing
	return &multipart.FileHeader{
		Filename: filename, // ✅ Giữ nguyên filename như input
		Size:     int64(len(data)),
		Header:   make(map[string][]string),
	}
}

// Test 1: Valid JPEG
func TestValidateAndSaveFile_ValidJPEG(t *testing.T) {
	uploadDir := t.TempDir()
	jpegData := createFakeJPEG(t)
	fileHeader := createFileHeader("test.jpg", jpegData)

	filename, err := ValidateAndSaveFile(fileHeader, uploadDir)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if filename == "" {
		t.Error("Expected filename, got empty string")
	}

	savedPath := filepath.Join(uploadDir, filename)
	if _, err := os.Stat(savedPath); os.IsNotExist(err) {
		t.Error("File was not saved")
	}

	t.Logf("✅ Valid JPEG test passed. File saved: %s", filename)
}

// Test 2: Valid PNG
func TestValidateAndSaveFile_ValidPNG(t *testing.T) {
	uploadDir := t.TempDir()
	pngData := createFakePNG(t)
	fileHeader := createFileHeader("test.png", pngData)

	filename, err := ValidateAndSaveFile(fileHeader, uploadDir)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	t.Logf("✅ Valid PNG test passed. File saved: %s", filename)
}

// Test 3: Invalid extension
func TestValidateAndSaveFile_InvalidExtension(t *testing.T) {
	uploadDir := t.TempDir()
	jpegData := createFakeJPEG(t)
	fileHeader := createFileHeader("test.exe", jpegData)

	_, err := ValidateAndSaveFile(fileHeader, uploadDir)

	if err == nil {
		t.Error("Expected error for .exe extension, got nil")
	}

	t.Logf("✅ Invalid extension blocked: %v", err)
}

// Test 4: File too large
func TestValidateAndSaveFile_FileTooLarge(t *testing.T) {
	uploadDir := t.TempDir()
	largeData := make([]byte, 6<<20) // 6MB

	// Tạo FileHeader thủ công với size lớn
	fileHeader := &multipart.FileHeader{
		Filename: "large.jpg",
		Size:     int64(len(largeData)),
	}

	_, err := ValidateAndSaveFile(fileHeader, uploadDir)

	if err == nil {
		t.Error("Expected error for large file, got nil")
	}

	t.Logf("✅ Large file blocked: %v", err)
}

// Test 5: Path traversal
func TestValidateAndSaveFile_PathTraversal(t *testing.T) {
	uploadDir := t.TempDir()
	maliciousNames := []string{
		"../../etc/passwd.jpg",
		"..\\..\\windows\\system32.jpg",
		"../../../secret.jpg",
		"..\\..\\..\\config.jpg",
	}

	for _, name := range maliciousNames {

		// Dùng helper để force exact filename
		fileHeader := &multipart.FileHeader{
			Filename: name,
			Size:     100,
		}
		_, err := ValidateAndSaveFile(fileHeader, uploadDir)

		if err == nil {
			t.Errorf("❌ Expected error for path traversal: %s", name)
		}
		t.Logf("✅ Blocked: %s → %v", name, err)
	}
}

// Test 6: Invalid filename characters
func TestValidateAndSaveFile_InvalidFilename(t *testing.T) {
	uploadDir := t.TempDir()

	invalidNames := []string{
		"test<script>.jpg",
		"test|file.jpg",
		"test?.jpg",
		`test"file.jpg`,
		"test:file.jpg",
		"test*file.jpg",
	}

	for _, name := range invalidNames {
		fileHeader := &multipart.FileHeader{
			Filename: name,
			Size:     100,
		}
		_, err := ValidateAndSaveFile(fileHeader, uploadDir)

		if err == nil {
			t.Errorf("Expected error for invalid filename: %s", name)
		}
		t.Logf("✅ Blocked: %s → %v", name, err)
	}
}

// Test 7: Fake JPEG (just random bytes)
func TestValidateAndSaveFile_FakeJPEG(t *testing.T) {
	uploadDir := t.TempDir()
	fakeData := []byte{0x4D, 0x5A, 0x90, 0x00} // MZ header
	fileHeader := createFileHeader("fake.jpg", fakeData)

	_, err := ValidateAndSaveFile(fileHeader, uploadDir)

	if err == nil {
		t.Error("Expected error for fake JPEG, got nil")
	}

	t.Logf("✅ Fake JPEG blocked: %v", err)
}

// Test 8: Corrupted image
func TestValidateAndSaveFile_CorruptedImage(t *testing.T) {
	uploadDir := t.TempDir()

	// JPEG magic bytes + garbage
	corrupted := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	corrupted = append(corrupted, make([]byte, 100)...)

	fileHeader := createFileHeader("corrupted.jpg", corrupted)

	_, err := ValidateAndSaveFile(fileHeader, uploadDir)

	if err == nil {
		t.Error("Expected error for corrupted image, got nil")
	}

	t.Logf("✅ Corrupted image blocked: %v", err)
}
