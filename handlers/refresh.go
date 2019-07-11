package handlers

import (
	"net/http"
	"time"
)

type refreshHandler struct {
	limiter chan<- time.Time
}

func (h *refreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.limiter <- time.Now()
	w.WriteHeader(http.StatusOK)
	return
}

func RefreshHandler(limiter chan<- time.Time) http.Handler {
	return &refreshHandler{limiter}
}
