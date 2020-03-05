package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/middleware"
	"net/http"

	htmux "github.com/dimfeld/httptreemux"
)

func FeedEvents(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if !middleware.CORS(w, r) {
		return
	}
	events, err := db.GetAllEvents(db.ConnectToDB())
	if err != nil {
		middleware.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	output, err := json.Marshal(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(output)
}

func EventHandler(router *htmux.TreeMux) {
	router.GET("/api/events/feed", FeedEvents)
}
