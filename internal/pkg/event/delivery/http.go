package delivery

import (
	"failless/internal/pkg/event"
	"failless/internal/pkg/event/usecase"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"log"
	"net/http"

	json "github.com/mailru/easyjson"
)

type eventDelivery struct {
	UseCase event.UseCase
}

func GetDelivery() event.Delivery {
	return &eventDelivery{
		UseCase: usecase.GetUseCase(),
	}
}

// Get events limited by number strings with offset from JSON (POST parameter)
// UserCount have to be set in the /configs/*/settings.go file using global variable
// UseCaseConf
func (ed *eventDelivery) GetEventsByKeyWords(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var searchRequest models.EventRequest
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}
	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	var events models.EventList
	if code, err := ed.UseCase.InitEventsByKeyWords(&events, searchRequest.Query, searchRequest.Page); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func (ed *eventDelivery) GetSearchEvents(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var searchRequest models.EventRequest
	err := json.UnmarshalFromReader(r.Body, &searchRequest)
	if err != nil {
		network.GenErrorCode(w, r, "Error within parse json", http.StatusBadRequest)
		return
	}

	log.Println(searchRequest)

	if searchRequest.Page < 1 {
		searchRequest.Page = 1
	}

	var events models.MidAndBigEventList
	if code, err := ed.UseCase.SearchEventsByUserPreferences(&events, &searchRequest); err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}
	network.Jsonify(w, events, http.StatusOK)
}

func (ed *eventDelivery) GetSmallEvents(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}


	events, err := ed.UseCase.GetSmallEventsByUID(int64(uid))
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	network.Jsonify(w, events, http.StatusOK)
}

func (ed *eventDelivery) CreateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	// Get form
	r.Header.Set("Content-Type", "application/json")
	var eventForm forms.SmallEventForm
	err := json.UnmarshalFromReader(r.Body, &eventForm)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	// Check form
	if !eventForm.Validate() {
		network.GenErrorCode(w, r, forms.MessageEventValidationFailed, http.StatusBadRequest)
		return
	}

	// Upload pics
	for iii := 0; iii < len(eventForm.Photos); iii++ {
		if eventForm.Photos[iii].ImgBase64 == "" ||
			!images.ValidateImage(&eventForm.Photos[iii], images.Events) {
			network.GenErrorCode(w, r, images.MessageImageValidationFailed, http.StatusBadRequest)
			return
		}
	}

	event, err := ed.UseCase.CreateSmallEvent(&eventForm)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	network.Jsonify(w, event, http.StatusOK)
}

func (ed *eventDelivery) UpdateSmallEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var event models.SmallEvent
	err := json.UnmarshalFromReader(r.Body, &event)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	code, err := ed.UseCase.UpdateSmallEvent(&event)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), code)
		return
	}

	network.Jsonify(w, event, code)
}

func (ed *eventDelivery) DeleteSmallEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	eid := network.GetEIdFromRequest(w, r, ps)
	if eid < 0 {
		return
	}

	message := ed.UseCase.DeleteSmallEvent(uid, eid)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, message, message.Status)
}

func (ed *eventDelivery) CreateMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	uid := security.CheckCredentials(w, r)
	if uid < 0 {
		return
	}

	r.Header.Set("Content-Type", "application/json")
	var midEventForm forms.MidEventForm
	err := json.UnmarshalFromReader(r.Body, &midEventForm)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}
	if !midEventForm.Validate() {
		network.GenErrorCode(w, r, forms.MessageEventValidationFailed, http.StatusBadRequest)
		return
	}

	// Upload pics
	for iii := 0; iii < len(midEventForm.Photos); iii++ {
		if midEventForm.Photos[iii].ImgBase64 == "" ||
			!images.ValidateImage(&midEventForm.Photos[iii], images.Events) {
			network.GenErrorCode(w, r, "Image validation failed. Check server logs", http.StatusBadRequest)
			return
		}
	}

	midEvent, message := ed.UseCase.CreateMidEvent(&midEventForm)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, midEvent, message.Status)
}

func (ed *eventDelivery) GetMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) UpdateMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) DeleteMiddleEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) JoinMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if uid := security.CheckCredentials(w, r); uid < 0 {
		return
	}

	if eid := network.GetEIdFromRequest(w, r, ps); eid < 0 {
		return
	}

	var subscription models.EventFollow
	err := json.UnmarshalFromReader(r.Body, &subscription)
	if err != nil {
		log.Println(err)
		network.GenErrorCode(w, r, network.MessageErrorParseJSON, http.StatusBadRequest)
		return
	}

	message := ed.UseCase.JoinMidEvent(&subscription)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, message, message.Status)
}

func (ed *eventDelivery) LeaveMiddleEvent(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	if uid := security.CheckCredentials(w, r); uid < 0 {
		return
	}

	if eid := network.GetEIdFromRequest(w, r, ps); eid < 0 {
		return
	}

	var subscription models.EventFollow
	err := json.UnmarshalFromReader(r.Body, &subscription)
	if err != nil {
		network.GenErrorCode(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	message := ed.UseCase.LeaveMidEvent(&subscription)
	if message.Message != "" {
		network.GenErrorCode(w, r, message.Message, message.Status)
		return
	}

	network.Jsonify(w, message, message.Status)
}

func (ed *eventDelivery) CreateBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) GetBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) UpdateBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) DeleteBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) AddVisitorForBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (ed *eventDelivery) RemoveVisitorForBigEvent(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	network.GenErrorCode(w, r, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
