package delivery

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/network"
	"net/http"
)

func FeedTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	//if !middleware.CORS(w, r) {
	//	return
	//}
	tags, err := db.GetAllTags(db.ConnectToDB())
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	network.Jsonify(w, tags, http.StatusOK)
	//output, err := json.Marshal(tags)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//w.Header().Set("content-type", "application/json")
	//_, _ = w.Write(output)
}

//func TagHandler(router *htmux.TreeMux) {
//	router.GET("/api/tags/feed", FeedTags)
//}
