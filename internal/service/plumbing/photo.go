package plumbing

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"io"
	"log/slog"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
)

func (rp *Plumping) GetImagePath(ctx context.Context, imageName string) (string, error) {
	imagePathStr := filepath.Join(rp.baseDir, "media", "images", imageName)

	imagePath, err := url.QueryUnescape(imagePathStr)
	if err != nil {
		fmt.Println("Error decoding query:", err)
	}
	if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("image not found: %s", imageName)
	}

	return imagePath, nil
}

func (pr *Plumping) UploadPhotos(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	const op = "service.UploadPhotos"
	var uploadedFiles []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			pr.log.Error("Failed to open file", slog.String("file", fileHeader.Filename), sl.Err(err))
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		defer file.Close()

		filePath := fmt.Sprintf("%s/media/images/%s", pr.baseDir, fileHeader.Filename)

		outFile, err := os.Create(filePath)
		if err != nil {
			pr.log.Error("Failed to create file", slog.String("file", filePath), sl.Err(err))
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			pr.log.Error("Failed to copy file data", slog.String("file", filePath), sl.Err(err))
			return nil, fmt.Errorf("%s, %w", op, err)
		}

		uploadedFiles = append(uploadedFiles, filePath)
	}

	return uploadedFiles, nil
}
