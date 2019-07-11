package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teambrookie/mediarss/db"
)

type mediaHandler struct {
	store db.MediaStore
}

func (h *mediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	episodes, err := h.store.GetCollection(db.FOUND)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(episodes)
	return
}

func MediaHandler(store db.MediaStore) http.Handler {
	return &mediaHandler{
		store: store,
	}
}
