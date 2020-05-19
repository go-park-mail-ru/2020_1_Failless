package usecase

import (
	"errors"
	"failless/internal/pkg/models"
	"failless/internal/pkg/tag"
	"failless/internal/pkg/tag/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func getTestUseCase(mockRep *mocks.MockRepository) tag.UseCase {
	return &tagUseCase{
		Rep: mockRep,
	}
}

func TestTagUseCase_InitEventsByTime_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	tc := getTestUseCase(mockRep)

	mockRep.EXPECT().GetAllTags().Return(nil, nil)

	status, err := tc.InitEventsByTime(new(models.TagList))

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, nil, err)
}

func TestTagUseCase_InitEventsByTime_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	tc := getTestUseCase(mockRep)

	correctError := errors.New("internal tag repo error")
	mockRep.EXPECT().GetAllTags().Return(nil, correctError)

	status, err := tc.InitEventsByTime(new(models.TagList))

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, correctError, err)
}
