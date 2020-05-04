package usecase

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func getTestUseCase(mockRep *mocks.MockRepository) chat.UseCase {
	return &chatUseCase{
		Rep: mockRep,
	}
}

func TestChatUseCase_Notify(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRep := mocks.NewMockRepository(mockCtrl)
	cc := getTestUseCase(mockRep)

	
}
