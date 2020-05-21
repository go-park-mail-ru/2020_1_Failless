package vote

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks failless/internal/pkg/vote Repository

type Repository interface {
	AddUserVote(uid int, id int, value int8) error
	CheckMatching(uid, id int) (bool, error)
}
