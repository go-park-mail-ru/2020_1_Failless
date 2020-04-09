package usecase

import (
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
		rep: mockRep,
	}
}

func TestTagUseCase_InitEventsByTime(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	vc := getTestUseCase(mockRep)

	mockRep.EXPECT().GetAllTags()

	status, err := vc.InitEventsByTime(new([]models.Tag))

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, nil, err)
}
