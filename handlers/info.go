package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teambrookie/mediarss/db"
)

type infoHandler struct {
	store db.MediaStore
}

func (h *infoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	if filename == "" {
		http.Error(w, "filename must be set in query params", http.StatusNotAcceptable)
		return
	}
	mediaInfo, err := h.store.GetMediaInfo(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mediaInfo)
	return

}

func InfoHandler(store db.MediaStore) http.Handler {
	return &infoHandler{store}
}
