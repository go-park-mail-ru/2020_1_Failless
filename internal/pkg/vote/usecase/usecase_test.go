package usecase

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestGetUseCase(t *testing.T) {

}

func TestVoteUseCase_VoteEvent(t *testing.T) {

}

func TestVoteUseCase_ValidateValue(t *testing.T) {
	vc := new(voteUseCase)

	assert.Equal(t, 1, vc.ValidateValue(1))
	assert.Equal(t, 1, vc.ValidateValue(math.MaxInt8))
	assert.Equal(t, 0, vc.ValidateValue(0))
	assert.Equal(t, -1, vc.ValidateValue(0.9))
	assert.Equal(t, -1, vc.ValidateValue(-math.MaxInt8))
}

func TestVoteUseCase_GetEventFollowers(t *testing.T) {

}
