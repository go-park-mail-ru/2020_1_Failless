package delivery

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/tag/usecase"
	"net/http"
)

func FeedTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uc := usecase.GetUseCase()
	var tags models.TagList
	if code, err := uc.InitEventsByTime(&tags); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, tags, http.StatusOK)
}
