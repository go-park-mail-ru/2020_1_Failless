package event

import "net/http"

type Delivery interface {
	// Search
	GetSearchEvents(w http.ResponseWriter, r *http.Request, _ map[string]string)

	// Small events
	GetSmallEvents(w http.ResponseWriter, r *http.Request, _ map[string]string)
	CreateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	UpdateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	DeleteSmallEvent(w http.ResponseWriter, r *http.Request, ps map[string]string)

	// Middle events
	CreateMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	GetMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	UpdateMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	DeleteMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	JoinMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string)
	LeaveMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string)

	// Big events
	CreateBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	GetBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	UpdateBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	DeleteBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	AddVisitorForBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
	RemoveVisitorForBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string)
}
