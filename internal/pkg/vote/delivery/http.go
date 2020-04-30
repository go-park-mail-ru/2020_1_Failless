package delivery

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/vote/usecase"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////// MANAGE ///////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////
func VoteUser(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("vote for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var vote models.Vote
	err := json.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	if uid != vote.Uid {
		network.GenErrorCode(w, r, "uid in the body is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.VoteUser(vote)
	network.Jsonify(w, message, message.Status)
}

func VoteEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("vote for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	var vote models.Vote
	err := json.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	if uid != vote.Uid {
		network.GenErrorCode(w, r, "uid in the body is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	message := uc.VoteEvent(vote)
	network.Jsonify(w, message, message.Status)
}

func EventFollowers(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	log.Println("follow for event")
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	id := network.GetIdFromRequest(w, r, ps)
	if id < 0 {
		network.GenErrorCode(w, r, "url id is incorrect", http.StatusBadRequest)
		return
	}

	uc := usecase.GetUseCase()
	followers, err := uc.GetEventFollowers(int(id))
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	network.Jsonify(w, followers, http.StatusOK)
}
