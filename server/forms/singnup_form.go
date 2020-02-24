package forms

import "regexp"

type SignForm struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

const (
	MinLen   = 6
	MinSym   = 4
	MinDigit = 2
)

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
	return true
}

func (s *SignForm) ValidateEmail() bool {
	r, _ := regexp.Compile("[a-zA-Z0-9.]+@[a-zA-Z0-9]+[.]{1}[a-z]{2,10}")
	return r.MatchString(s.Email)
}

func (s *SignForm) ValidatePhone() bool {
	digitCounter := 0
	for _, sym := range s.Phone {
		if '0' <= sym && sym <= '9' {
			digitCounter++
		}
		return false
	}
	if !(5 < digitCounter && digitCounter < 15) {
		return false
	}
	return true
}

func (s *SignForm) Validate() bool {
	return s.ValidateEmail() && s.ValidatePassword() && s.ValidatePhone()
}
