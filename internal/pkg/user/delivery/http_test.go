package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/images"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	testSignForm = forms.SignForm{
		Uid:      security.TestUser.Uid,
		Name:     security.TestUser.Name,
		Phone:    security.TestUser.Phone,
		Email:    security.TestUser.Email,
		Password: "qwerty1234",
	}
	testInvalidUidType = map[string]interface{}{
		"uid": strconv.Itoa(1),			// Invalid type
	}
	useCaseError = errors.New("error in usecase")
	testMessageUseCaseError = models.WorkMessage{
		Request: nil,
		Message: useCaseError.Error(),
		Status:  http.StatusInternalServerError,
	}
	testAbout = models.UserAbout{About:"about"}
	testUserTags = models.UserTags{Tags:[]int{1,2}}
	testUserPhotos = forms.EImageList{{ImgBase64:""}}
	testUserPhotosInvalid = forms.EImageList{{ImgBase64:"kek", ImgName:""}}
	testGeneralForm = forms.GeneralForm{
		SignForm: testSignForm,
		Tags:     nil,
		Avatar:   forms.EImage{},
		Photos:   nil,
		Gender:   0,
		About:    "",
		Rating:   0,
		Location: models.LocationPoint{},
		Birthday: time.Time{},
	}
	testUserRequest = models.UserRequest{
		Uid:      security.TestUser.Uid,
		Page:     1,
		Limit:    10,
		Query:    "kek",
		Tags:     nil,
		Location: models.LocationPoint{},
		MinAge:   18,
		MaxAge:   100,
		Men:      true,
		Women:    true,
	}
)

func getTestDelivery(mockUC *mocks.MockUseCase) user.Delivery {
	return &userDelivery{UseCase:mockUC}
}

func TestUpdProfileGeneral_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("PUT", "/api/profile/:id/general", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.UpdProfileGeneral(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_UpdProfileGeneral_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("PUT", "/api/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	ud.UpdProfileGeneral(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_UpdProfileGeneral_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSignForm)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserBase(&testSignForm).Return(testMessageUseCaseError.Status, useCaseError)
	ud.UpdProfileGeneral(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_UpdProfileGeneral_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testSignForm)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/general", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserBase(&testSignForm)
	ud.UpdProfileGeneral(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_UpdUserAbout_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.UpdUserAbout(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_UpdUserAbout_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	testInvalidAboutType := map[string]interface{}{
		"about": 1,			// Invalid type
	}
	body, _ := json.Marshal(testInvalidAboutType)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	ud.UpdUserAbout(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_UpdUserAbout_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testAbout)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserAbout(security.TestUser.Uid, testAbout.About).Return(testMessageUseCaseError)
	ud.UpdUserAbout(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_UpdUserAbout_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testAbout)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserAbout(security.TestUser.Uid, testAbout.About)
	ud.UpdUserAbout(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_UpdUserTags_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.UpdUserTags(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_UpdUserTags_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	testInvalidUserTagsType := map[string]interface{}{
		"tags": 1,			// Invalid type
	}
	body, _ := json.Marshal(testInvalidUserTagsType)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	ud.UpdUserTags(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_UpdUserTags_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserTags)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserTags(security.TestUser.Uid, testUserTags.Tags).Return(testMessageUseCaseError)
	ud.UpdUserTags(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_UpdUserTags_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserTags)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/about", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserTags(security.TestUser.Uid, testUserTags.Tags)
	ud.UpdUserTags(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_UpdUserPhotos_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/photos", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.UpdUserPhotos(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_UpdUserPhotos_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	testInvalidUserPhotosType := map[string]interface{}{
		"photos": 1,			// Invalid type
	}
	body, _ := json.Marshal(testInvalidUserPhotosType)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/photos", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	ud.UpdUserPhotos(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_UpdUserPhotos_IncorrectPhotos(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotosInvalid)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/photos", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	ud.UpdUserPhotos(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusNotFound, msg.Status)
	assert.Equal(t, images.MessageImageValidationFailed, msg.Message)
}

func TestUserDelivery_UpdUserPhotos_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/photos", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserPhotos(security.TestUser.Uid, &testUserPhotos).Return(testMessageUseCaseError)
	ud.UpdUserPhotos(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_UpdUserPhotos_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("PUT", "/api/srv/profile/:id/meta/photos", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().UpdateUserPhotos(security.TestUser.Uid, &testUserPhotos)
	ud.UpdUserPhotos(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetProfilePage_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.GetProfilePage(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_GetProfilePage_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("GET", "/api/srv/profile/:id", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	testForm := forms.GeneralForm{}
	testForm.Uid = testSignForm.Uid
	mockUC.EXPECT().GetUserInfo(&testForm).Return(testMessageUseCaseError.Status, useCaseError)
	ud.GetProfilePage(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_GetProfilePage_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("GET", "/api/srv/profile/:id", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	testForm := forms.GeneralForm{}
	testForm.Uid = testSignForm.Uid
	mockUC.EXPECT().GetUserInfo(&testForm)
	ud.GetProfilePage(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetSmallEventsForUser_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id/small-events", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.GetSmallEventsForUser(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_GetSmallEventsForUser_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id/small-events", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().GetSmallEventsForUser(new(models.SmallEventList), testSignForm.Uid).Return(testMessageUseCaseError)
	ud.GetSmallEventsForUser(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_GetSmallEventsForUser_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("GET", "/api/srv/profile/:id/small-events", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().GetSmallEventsForUser(new(models.SmallEventList), testSignForm.Uid)
	ud.GetSmallEventsForUser(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetSmallAndMidEventsForUser_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id/own-events", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.GetSmallAndMidEventsForUser(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_GetSmallAndMidEventsForUser_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id/own-events", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().GetUserOwnEvents(new(models.OwnEventsList), testSignForm.Uid).Return(testMessageUseCaseError)
	ud.GetSmallAndMidEventsForUser(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_GetSmallAndMidEventsForUser_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("GET", "/api/srv/profile/:id/own-events", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().GetUserOwnEvents(new(models.OwnEventsList), testSignForm.Uid)
	ud.GetSmallAndMidEventsForUser(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetProfileSubscriptions_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id/subscriptions", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": "kek"}

	ud.GetProfileSubscriptions(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageInvalidID, msg.Message)
}

func TestUserDelivery_GetProfileSubscriptions_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/profile/:id/subscriptions", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().GetUserSubscriptions(new(models.MidAndBigEventList), testSignForm.Uid).Return(testMessageUseCaseError)
	ud.GetProfileSubscriptions(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_GetProfileSubscriptions_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("GET", "/api/srv/profile/:id/subscriptions", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)
	ps := map[string]string{"id": strconv.Itoa(testSignForm.Uid)}

	mockUC.EXPECT().GetUserSubscriptions(new(models.MidAndBigEventList), testSignForm.Uid)
	ud.GetProfileSubscriptions(rr, req.WithContext(ctx), ps)

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetUserInfo_Incorrect(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUC)

	req, err := http.NewRequest("GET", "/api/srv/getuser", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ud.GetUserInfo(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestUserDelivery_GetUserInfo_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ud := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	body, _ := json.Marshal(testUserPhotos)
	req, err := http.NewRequest("GET", "/api/srv/getuser", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ud.GetUserInfo(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_SignIn_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ud := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/srv/signin", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ud.SignIn(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_SignIn_InvalidForm(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ud := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	invForm := testSignForm
	invForm.Phone = "0000000000000000"
	invForm.Email = ""
	body, _ := json.Marshal(invForm)
	req, err := http.NewRequest("POST", "/api/srv/signin", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ud.SignIn(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageValidationFailed, msg.Message)
}

func TestUserDelivery_Logout_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ud := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("GET", "/api/srv/logout", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ud.Logout(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusOK, msg.Status)
	assert.Equal(t, network.MessageSuccessfulLogout, msg.Message)
}

func TestUserDelivery_SignUp_CorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ed := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("POST", "/api/srv/signup", nil)
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.SignUp(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_SignUp_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ed := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/srv/signup", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()

	ed.SignUp(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_SignUp_InvalidForm(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ud := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	invForm := testSignForm
	invForm.Phone = "0000000000000000"
	body, _ := json.Marshal(invForm)
	req, err := http.NewRequest("POST", "/api/srv/signup", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ud.SignUp(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageValidationFailed, msg.Message)
}

func TestUserDelivery_SignUp_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUc := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUc)

	body, _ := json.Marshal(testSignForm)
	req, err := http.NewRequest("POST", "/api/srv/signup", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	usr := models.User{
		Uid:      -1,
		Phone:    testSignForm.Phone,
		Email:    testSignForm.Email,
	}
	mockUc.EXPECT().FillFormIfExist(&usr).Return(testMessageUseCaseError.Status, useCaseError)
	ud.SignUp(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_SignUp_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUc := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUc)

	body, _ := json.Marshal(testSignForm)
	req, err := http.NewRequest("POST", "/api/srv/signup", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	usr := models.User{
		Uid:      -1,
		Phone:    testSignForm.Phone,
		Email:    testSignForm.Email,
	}
	mockUc.EXPECT().FillFormIfExist(&usr).Return(http.StatusOK, nil)
	mockUc.EXPECT().RegisterNewUser(&testSignForm).Return(useCaseError)
	ud.SignUp(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusInternalServerError, msg.Status)
}

func TestUserDelivery_SignUp_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUc := mocks.NewMockUseCase(mockCtrl)
	ud := getTestDelivery(mockUc)

	body, _ := json.Marshal(testSignForm)
	req, err := http.NewRequest("POST", "/api/srv/signup", bytes.NewReader(body))
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	usr := models.User{
		Uid:      -1,
		Phone:    testSignForm.Phone,
		Email:    testSignForm.Email,
	}
	mockUc.EXPECT().FillFormIfExist(&usr).Return(http.StatusOK, nil)
	mockUc.EXPECT().RegisterNewUser(&testSignForm).Return(nil)
	ud.SignUp(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetUsersFeed_IncorrectUid(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ud := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	req, err := http.NewRequest("POST", "/api/srv/users/feed", nil)
	if err != nil {
		t.Fail()
		return
	}
	rr := httptest.NewRecorder()

	ud.GetUsersFeed(rr, req, map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusUnauthorized, msg.Status)
	assert.Equal(t, network.MessageErrorAuthRequired, msg.Message)
}

func TestUserDelivery_GetUsersFeed_IncorrectBody(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ed := getTestDelivery(mocks.NewMockUseCase(mockCtrl))

	body, _ := json.Marshal(testInvalidUidType)
	req, err := http.NewRequest("POST", "/api/srv/users/feed", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	ed.GetUsersFeed(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, http.StatusBadRequest, msg.Status)
	assert.Equal(t, network.MessageErrorParseJSON, msg.Message)
}

func TestUserDelivery_GetUsersFeed_IncorrectPage(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	searchReq := testUserRequest
	searchReq.Page = 0
	body, _ := json.Marshal(searchReq)
	req, err := http.NewRequest("POST", "/api/srv/users/feed", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest)
	mockUC.EXPECT().GetFeedResultsFor(security.TestUser.Uid, new([]models.UserGeneral))
	ed.GetUsersFeed(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}

func TestUserDelivery_GetUsersFeed_Incorrect1(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserRequest)
	req, err := http.NewRequest("POST", "/api/srv/users/feed", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest).Return(testMessageUseCaseError.Status, useCaseError)
	ed.GetUsersFeed(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError.Status, msg.Status)
	assert.Equal(t, testMessageUseCaseError.Message, msg.Message)
}

func TestUserDelivery_GetUsersFeed_Incorrect2(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserRequest)
	req, err := http.NewRequest("POST", "/api/srv/users/feed", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest)
	mockUC.EXPECT().GetFeedResultsFor(security.TestUser.Uid, new([]models.UserGeneral)).Return(nil, testMessageUseCaseError)
	ed.GetUsersFeed(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, testMessageUseCaseError, msg)
}

func TestUserDelivery_GetUsersFeed_Correct(t *testing.T) {
	// Create mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUC := mocks.NewMockUseCase(mockCtrl)
	ed := getTestDelivery(mockUC)

	body, _ := json.Marshal(testUserRequest)
	req, err := http.NewRequest("POST", "/api/srv/users/feed", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
		return
	}
	rr := httptest.NewRecorder()
	ctx := context.Background()
	ctx = context.WithValue(ctx, security.CtxUserKey, security.TestUser)

	mockUC.EXPECT().InitUsersByUserPreferences(new([]models.UserGeneral), &testUserRequest)
	mockUC.EXPECT().GetFeedResultsFor(security.TestUser.Uid, new([]models.UserGeneral))
	ed.GetUsersFeed(rr, req.WithContext(ctx), map[string]string{})

	msg, err := network.DecodeToMsg(rr.Body)
	if err != nil {
		t.Fail()
		return
	}

	assert.Equal(t, 0, msg.Status)
}
