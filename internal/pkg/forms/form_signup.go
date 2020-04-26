package forms

import (
	pb "failless/api/proto/auth"
	"failless/internal/pkg/models"
	"log"
	"regexp"
	"sync"
)

type SignForm struct {
	Uid      int    `json:"uid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password, omitempty"`
}

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
	log.Println("email invalid")
	s.Email = ""
	return false
}

func (s *SignForm) ValidatePhone() bool {
	digitCounter := 0
	log.Println(s.Phone)
	for _, sym := range s.Phone {
		if '0' <= sym && sym <= '9' {
			digitCounter++
		} else {
			return false
		}
	}
	log.Println(digitCounter)
	if 5 < digitCounter && digitCounter < 15 {
		log.Println("phone valid")
		return true
	}
	s.Phone = ""
	log.Println("phone invalid")
	return false
}

func (s *SignForm) Validate() bool {
	return s.ValidateEmail() && s.ValidatePassword() && s.ValidatePhone()
}

func (s *SignForm) FillFromModel(user *models.User) {
	s.Password = ""
	s.Name = user.Name
	s.Email = user.Email
	s.Phone = user.Phone
	s.Uid = user.Uid
}

func (s *SignForm) FillFromAuthReply(user *pb.AuthorizeReply) {
	s.Password = ""
	s.Name = user.Cred.Name
	s.Email = user.Cred.Email
	s.Phone = user.Cred.Phone
	s.Uid = int(user.Cred.Uid)
}
