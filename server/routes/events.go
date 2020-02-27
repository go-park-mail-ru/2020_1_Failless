package routes

import (
	"encoding/json"
	"failless/db"
	htmux "github.com/dimfeld/httptreemux"
	"net/http"
)

func FeedEvents(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	CORS(w, r)
	events, err := db.GetAllEvents(db.ConnectToDB())
	if err != nil {
		GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
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

