package user

import "net/http"

type Delivery interface {
	// Get
	GetProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string)
	GetUserInfo(w http.ResponseWriter, r *http.Request, _ map[string]string)

	// Authorization
	SignIn(w http.ResponseWriter, r *http.Request, _ map[string]string)
	Logout(w http.ResponseWriter, r *http.Request, _ map[string]string)
	SignUp(w http.ResponseWriter, r *http.Request, _ map[string]string)

	// Update
	UpdProfileGeneral(w http.ResponseWriter, r *http.Request, ps map[string]string)
	UpdUserAbout(w http.ResponseWriter, r *http.Request, ps map[string]string)
	UpdUserTags(w http.ResponseWriter, r *http.Request, ps map[string]string)
	UpdProfilePage(w http.ResponseWriter, r *http.Request, ps map[string]string)
	UpdUserPhotos(w http.ResponseWriter, r *http.Request, ps map[string]string)

	// Events
	GetSmallEventsForUser(w http.ResponseWriter, r *http.Request, ps map[string]string)
	GetSmallAndMidEventsForUser(w http.ResponseWriter, r *http.Request, ps map[string]string)
	GetProfileSubscriptions(w http.ResponseWriter, r *http.Request, ps map[string]string)

	// Feed
	GetUsersFeed(w http.ResponseWriter, r *http.Request, _ map[string]string)
}
