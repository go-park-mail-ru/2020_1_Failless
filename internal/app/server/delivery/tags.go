package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/network"
	"net/http"

	htmux "github.com/dimfeld/httptreemux"
)

func FeedTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	if !network.CORS(w, r) {
		return
	}
	tags, err := db.GetAllTags(db.ConnectToDB())
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
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
