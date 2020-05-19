package tag

import (
	"net/http"
)

type Delivery interface {
	FeedTags(w http.ResponseWriter, r *http.Request, _ map[string]string)
}