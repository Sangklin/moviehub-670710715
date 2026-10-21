package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"moviehub/mw"
	"moviehub/store"
)

func (h *Handler) PutVote(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	var req struct {
		Score int `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
		return
	}
	if req.Score < 1 || req.Score > 10 {
		writeError(w, http.StatusBadRequest, "BAD_SCORE", "score ต้องเป็นจำนวนเต็ม 1 ถึง 10")
		return
	}
	who := mw.WhoFrom(r.Context())
	err := h.Store.PutVote(r.Context(), who.ID, movieID, req.Score)
	if errors.Is(err, store.ErrNoMovie) {
		writeError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "ไม่พบหนังเรื่องนี้")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"movieId": movieID, "score": req.Score})
}

func (h *Handler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "id")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	who := mw.WhoFrom(r.Context())
	if err := h.Store.DeleteVote(r.Context(), who.ID, movieID); err != nil {
		internal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
