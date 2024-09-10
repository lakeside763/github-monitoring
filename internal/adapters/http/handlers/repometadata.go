package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/lakeside763/github-repo/internal/core/ports/interfaces"
	"github.com/lakeside763/github-repo/pkg/utils"
)

type RepometadataHandler struct {
	Repository interfaces.Repometadata
}

func NewRepometadataHandler(repo interfaces.Repometadata) *RepometadataHandler {
	return &RepometadataHandler{Repository: repo}
}

func (h *RepometadataHandler) CreateRepometadata(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var repo models.Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	newRepo, err := h.Repository.Create(&repo);
	if err != nil {
		http.Error(w, "Failed to create repo metadata", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, http.StatusOK, newRepo)
}