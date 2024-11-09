package plumbing

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
	"os"
)

func (s *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Path[len("/media/image/"):]

	imagePath, err := s.service.GetImagePath(context.Background(), imageName)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found image"}, http.StatusNotFound)
		return
	}

	imgFile, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open image: %v", err), http.StatusInternalServerError)
		return
	}
	defer imgFile.Close()

	fileStat, err := imgFile.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to stat image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))

	http.ServeContent(w, r, imageName, fileStat.ModTime(), imgFile)
}
