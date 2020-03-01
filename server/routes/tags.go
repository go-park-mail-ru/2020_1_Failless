package routes

import (
	"encoding/json"
	"failless/db"
	htmux "github.com/dimfeld/httptreemux"
	"net/http"
)

func FeedTags(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if !CORS(w, r) {
		return
	}
	tags, err := db.GetAllTags(db.ConnectToDB())
	if err != nil {
		GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(output)
}

func TagHandler(router *htmux.TreeMux) {
	router.GET("/api/tags/feed", FeedTags)
}
