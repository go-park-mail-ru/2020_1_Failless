package delivery

import (
	"failless/internal/pkg/event/usecase"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"net/http"
)

func FeedEvents(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uc := usecase.GetUseCase()
	var events []models.Event
	if code, err := uc.InitEventsByTime(events); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}
