package delivery

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/network"
	"net/http"
)

func FeedTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	tags, err := db.GetAllTags(db.ConnectToDB())
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	network.Jsonify(w, tags, http.StatusOK)
}
