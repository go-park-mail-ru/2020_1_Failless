package forms

import (
	"failless/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"sync"
)

type SignForm struct {
	Name     string `json:"name, omitempty"`
	Phone    string `json:"phone, omitempty"`
	Email    string `json:"email, omitempty"`
	Password string `json:"password, omitempty"`
}

const (
	MinLen   = 6
	MinSym   = 4
	MinDigit = 2
)

var compileOnce = sync.Once{}
var regExpr *regexp.Regexp = nil

func (s *SignForm) ValidatePassword() bool {
	if len(s.Password) < MinLen {
		return false
	}
	digitCounter := 0
	symCounter := 0
	for _, elem := range s.Password {
		if '0' <= elem && elem <= '9' {
			digitCounter++
		} else {
			symCounter++
		}
	}
	if symCounter < MinSym || digitCounter < MinDigit {
		return false
	}
	log.Println("password valid")
	return true
}

func (s *SignForm) ValidateEmail() bool {
	compileOnce.Do(func() {
		regExpr, _ = regexp.Compile("[a-zA-Z0-9.]+@[a-zA-Z0-9]+[.]{1}[a-z]{2,10}")
	})
	if regExpr.MatchString(s.Email) {
		log.Println("email valid")
		return true
	}
	return false
}

func (s *SignForm) ValidatePhone() bool {
	digitCounter := 0
	for _, sym := range s.Phone {
		if '0' <= sym && sym <= '9' {
			digitCounter++
		} else {
			return false
		}
	}
	if !(5 < digitCounter && digitCounter < 15) {
		return false
	}
	log.Println("phone valid")
	return true
}

func (s *SignForm) Validate() bool {
	return s.ValidateEmail() && s.ValidatePassword() && s.ValidatePhone()
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func RegisterNewUser(user SignForm) error {
	bPass, err := EncryptPassword(user.Password)
	if err != nil {
		return err
	}

	dbUser := db.User{
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: bPass,
	}

	return db.AddNewUser(db.ConnectToDB(), &dbUser)
}

