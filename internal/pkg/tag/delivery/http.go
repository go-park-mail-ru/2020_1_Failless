package delivery

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/usecase"
	"net/http"
)

type tagDelivery struct {
	UseCase tag.UseCase
}

func GetDelivery() tag.Delivery {
	return &tagDelivery{
		UseCase: usecase.GetUseCase(),
	}
}

func (td *tagDelivery) FeedTags(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var tags models.TagList
	if code, err := td.UseCase.InitEventsByTime(&tags); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, tags, http.StatusOK)
}
