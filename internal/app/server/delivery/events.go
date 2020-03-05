package delivery

import (
	"encoding/json"
	"failless/internal/pkg/db"
	"failless/internal/pkg/network"
	"net/http"
)

func FeedEvents(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	//if !middleware.CORS(w, r) {
	//	return
	//}
	events, err := db.GetAllEvents(db.ConnectToDB())
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
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

//func EventHandler(router *htmux.TreeMux) {
//	router.GET("/api/events/feed", FeedEvents)
//}
