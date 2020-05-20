package delivery

import (
	"failless/internal/pkg/chat"
	"failless/internal/pkg/chat/mocks"
	"testing"
)

func getTestDelivery(mockUC *mocks.MockUseCase) chat.Delivery {
	return &chatDelivery{UseCase:mockUC}
}

func TestChatDelivery_GetMessages(t *testing.T) {

}
