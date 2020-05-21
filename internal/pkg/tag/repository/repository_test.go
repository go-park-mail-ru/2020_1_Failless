package repository

import (
	"errors"
	"failless/internal/pkg/db/mocks"
	"failless/internal/pkg/tag"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testDBError = errors.New("db error")
)

func getTestRep(mockDB *mocks.MockMyDBInterface) tag.Repository {
	return &sqlTagRepository{
		db: mockDB,
	}
}

func TestSqlTagRepository_GetAllTags_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDB := mocks.NewMockMyDBInterface(mockCtrl)
	tr := getTestRep(mockDB)

	mockDB.EXPECT().Query(QuerySelectTagById).Return(nil, testDBError)

	_, err := tr.GetAllTags()

	assert.Equal(t, testDBError, err)
}
