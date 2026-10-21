package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"unicode/utf8"

	"moviehub/mw"
	"moviehub/store"
)

func (h *Handler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	memberID, _ := mw.MemberID(r.Context())
	items, err := h.Store.GetWishlist(r.Context(), memberID)
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *Handler) PutWishlist(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "movieId")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	var req struct {
		Note *string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
		return
	}
	if req.Note != nil && utf8.RuneCountInString(*req.Note) > 200 {
		writeError(w, http.StatusBadRequest, "BAD_NOTE", "โน้ตต้องยาวไม่เกิน 200 ตัวอักษร")
		return
	}
	memberID, _ := mw.MemberID(r.Context())
	err := h.Store.PutWishlist(r.Context(), memberID, movieID, req.Note)
	if errors.Is(err, store.ErrNoMovie) {
		writeError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "ไม่พบหนังเรื่องนี้")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"movieId": movieID, "note": req.Note})
}

func (h *Handler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "movieId")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	memberID, _ := mw.MemberID(r.Context())
	if err := h.Store.DeleteWishlist(r.Context(), memberID, movieID); err != nil {
		internal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
