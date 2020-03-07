package delivery

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/network"
	"net/http"
)

func FeedEvents(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	events, err := db.GetAllEvents(db.ConnectToDB())
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}
