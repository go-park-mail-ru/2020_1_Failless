package delivery

import (
	"encoding/json"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/vote/usecase"
	"net/http"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////// MANAGE ///////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////
func VoteUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
}

func VoteEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var vote models.Vote
	err := decoder.Decode(&vote)
	if err != nil {
		network.GenErrorCode(w, r,"Error within parse json", http.StatusBadRequest)
		return
	}

	if uid != vote.Uid {
		network.GenErrorCode(w, r, "uid in the body is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.VoteEvent(vote)
	if message.Status != http.StatusOK {
		network.Jsonify(w, message, message.Status)
		return
	}

}

func FollowEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
}
