package images



import (
	"fmt"
	// "image"
	// "log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// ✅ FilterImages applies filters and fixes orientation.
func FilterImages(inputFilepath string) (string, error) {
	// ✅ Open image with auto orientation
	src, err := imaging.Open(inputFilepath, imaging.AutoOrientation(true))
	// src, err := OpenWithAutoOrientation(inputFilepath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}

	// ✅ Apply some enhancements
	filtered := imaging.AdjustContrast(src, 15)
	filtered = imaging.Sharpen(filtered, 1.5)
	filtered = imaging.CropCenter(filtered,200,200);
	// Optional resize
	if filtered.Bounds().Dx() > 1200 {
		filtered = imaging.Resize(filtered, 1200, 0, imaging.Lanczos)
	}

	uploadDir := "./uploads"
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	fileUUID := uuid.New().String()
	outputPath := filepath.Join(uploadDir, fmt.Sprintf("%s_filtered.jpg", fileUUID))

	err = imaging.Save(filtered, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	// log.Printf("✅ Saved filtered image: %s", outputPath)
	return outputPath, nil
}



