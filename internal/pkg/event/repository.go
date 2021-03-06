package event

import "failless/internal/pkg/models"

type Repository interface {
	GetNameByID(uid int) (string, error)
	GetValidTags() ([]models.Tag, error)
	JoinMidEvent(uid, eid int) (int, error)
	LeaveMidEvent(uid, eid int) (int, error)
	FollowBigEvent(uid, eid int) error
	UnfollowBigEvent(uid, eid int) error
	CreateSmallEvent(event *models.SmallEvent) error
	UpdateSmallEvent(event *models.SmallEvent) (int, error)
	DeleteSmallEvent(uid int, eid int64) error
	GetSmallEventsForUser(smallEvents *models.SmallEventList, uid int) (int, error)
	CreateMidEvent(event *models.MidEvent) error
	GetOwnMidEvents(midEvents *models.MidEventList, uid int) (int, error)
	GetAllMidEvents(midEvents *models.MidEventList, request *models.EventRequest) (int, error)
	GetSubscriptionMidEvents(midEvent *models.MidEventList, uid int) (int, error)
	GetMidEventsWithFollowed(midEvents *models.MidEventList, request *models.EventRequest) (int, error)
	GetOwnMidEventsWithAnotherUserFollowed(midEvents *models.MidEventList, admin, member int) (int, error)
	GetSubscriptionMidEventsWithAnotherUserFollowed(midEvents *models.MidEventList, uid, visitor int) (int, error)
}
