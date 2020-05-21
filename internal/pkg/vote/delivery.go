package vote

import (
	"net/http"
)

type Delivery interface {
	VoteUser(w http.ResponseWriter, r *http.Request, ps map[string]string)
	MatchPush(w http.ResponseWriter, r *http.Request, ps map[string]string)
}