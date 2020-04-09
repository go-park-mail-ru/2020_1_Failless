package delivery

import (
	"failless/internal/pkg/forms"
	"failless/internal/pkg/models"
	"failless/internal/pkg/security"
	"failless/internal/pkg/user"
	"failless/internal/pkg/user/mocks"
	"failless/internal/pkg/user/usecase"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSignUp2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUseCase := mocks.NewMockRepository(mockCtrl)
	form := forms.SignForm{
		Uid:      1,
		Name:     "MrTester",
		Phone:    "88005553535",
		Email:    "mrtester@mr.tester",
		Password: "qwerty12345",
	}

	callFirs := mockUseCase.EXPECT().GetUserByPhoneOrEmail(
		form.Phone,
		form.Email).Return(
		models.User{Uid: -1, Password: []byte{}}, nil).Times(1)

	bPass, _ := security.EncryptPassword(form.Password)

	mockUseCase.EXPECT().AddNewUser(&models.User{
		Name:     form.Name,
		Phone:    form.Phone,
		Email:    form.Email,
		Password: bPass,
	}).After(callFirs)
	var useCase user.UseCase = &usecase.UserUseCase{
		Rep: mockUseCase,
	}

	err := useCase.RegisterNewUser(&form)
	assert.Nil(t, err)
}
