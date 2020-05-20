package usecase

import (
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/mocks"
)

func getTestUseCase(mockRep *mocks.MockRepository) event.UseCase {
	return &eventUseCase{
		rep: mockRep,
	}
}


