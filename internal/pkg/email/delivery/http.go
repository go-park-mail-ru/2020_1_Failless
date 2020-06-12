package delivery

import (
	"failless/internal/pkg/email"
	"failless/internal/pkg/email/usecase"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	json "github.com/mailru/easyjson"
	"net/http"
)

type emailDelivery struct {
	UseCase email.UseCase
}

func GetDelivery() email.Delivery {
	return &emailDelivery{
		UseCase: usecase.GetUseCase(),
	}
}

func (ed *emailDelivery) SendReminder(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var receiver models.Email
	err := json.UnmarshalFromReader(r.Body, &receiver)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	if msg := ed.UseCase.SendReminder(&receiver); msg.Status >= 400 {
		network.GenErrorCode(w, r, msg.Message, msg.Status)
		return
	}
	network.Jsonify(w, receiver, http.StatusOK)
}