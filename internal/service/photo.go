package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (rp *Plumping) GetImagePath(ctx context.Context, imageName string) (string, error) {
	imagePath := filepath.Join(rp.baseDir, "media", "images", imageName)

	if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("image not found: %s", imageName)
	}

	return imagePath, nil
}
