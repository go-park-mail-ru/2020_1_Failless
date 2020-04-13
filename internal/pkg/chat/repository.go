package chat

type Repository interface {
	InsertDialogue(id1, id2 int) (int, error)
}