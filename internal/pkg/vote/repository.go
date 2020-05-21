package vote

type Repository interface {
	AddUserVote(uid int, id int, value int8) error
	CheckMatching(uid, id int) (bool, error)
}
