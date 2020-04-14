package chat

type UseCase interface {
	CreateDialogue(id1, id2 int) (int, error)
}
