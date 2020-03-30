package security

import (
	"failless/internal/pkg/middleware"
	"failless/internal/pkg/network"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ComparePasswords(hash []byte, p string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(p))
	return err == nil
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckCredentials(w http.ResponseWriter, r *http.Request) int {
	data := r.Context().Value(middleware.CtxUserKey)
	if data == nil {
		network.GenErrorCode(w, r, "auth required", http.StatusUnauthorized)
		return -1
	}

	cred := data.(middleware.UserClaims)
	if cred.Uid < 0 {
		network.GenErrorCode(w, r, "token uid is incorrect", http.StatusBadRequest)
		return -1
	}

	return cred.Uid
}

func CompareUidsFromURLAndToken(w http.ResponseWriter, r *http.Request, ps map[string]string) int {
	ctxUid := CheckCredentials(w, r)

	uid := 0
	if uid = network.GetIdFromRequest(w, r, &ps); uid < 0 {
		network.GenErrorCode(w, r, "url uid is incorrect", http.StatusBadRequest)
		return -1
	}

	if ctxUid != uid {
		network.GenErrorCode(w, r, "forbidden", http.StatusForbidden)
		return -1
	}

	return uid
}


// TODO: add to header CSRF-token when create request for up vote and down vote and all post|put|delete requests
// login | register pages not needed in csrf token
// and if we tape reload in the browser our SPA have to still work fine
// we get token without session, using separate method for this (get token and after this go to handler)
// we may save token to local storage in browser but not into page memory

