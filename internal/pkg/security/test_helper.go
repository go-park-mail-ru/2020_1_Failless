package security

import "errors"

var (
	TestUser = UserClaims{
		Uid:   1,
		Phone: "88005553535",
		Email: "mail@mail.ru",
		Name:  "mrTester",
	}
	InvalidTestUser = UserClaims{
		Uid:   -1,
		Phone: "00000000000000",
		Email: "",
		Name:  "",
	}
	claimsNotFoundError = errors.New("Claims not found\n")
	incorrectTokenUidError = errors.New("token uid is incorrect")
)
