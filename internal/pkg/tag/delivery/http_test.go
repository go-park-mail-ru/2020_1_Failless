package delivery

import (
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getTestDelivery(mockUC *mocks.MockUseCase) tag.Delivery {
	return &tagDelivery{UseCase:mockUC}
}

type TestCaseTags struct {
	Response   string
	StatusCode int
}

func TestTagDelivery_FeedTags_InternalError(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	td := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/users/:vote", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	var tags models.TagList
	mockUC.EXPECT().InitEventsByTime(&tags).Return(http.StatusInternalServerError, errors.New("internal tag error"))

	td.FeedTags(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
	assert.Equal(t, "internal tag error", msg.Message)
}

func TestTagDelivery_FeedTags_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	td := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/users/:vote", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	var ps map[string]string

	var tags models.TagList
	mockUC.EXPECT().InitEventsByTime(&tags)

	td.FeedTags(rr, req, ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}
