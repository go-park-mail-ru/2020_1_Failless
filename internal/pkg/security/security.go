package security

import (
	"crypto/rand"
	"failless/internal/pkg/network"
	"failless/internal/pkg/settings"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserClaims struct {
	Uid   int
	Phone string
	Email string
	Name  string
}

type UserKey string

// Context variable for pushing credentials through middleware to handlers
const CtxUserKey UserKey = "auth"

func ComparePasswords(hash []byte, p string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(p))
	return err == nil
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckCredentials(w http.ResponseWriter, r *http.Request) int {
	data := r.Context().Value(CtxUserKey)
	if data == nil {
		network.GenErrorCode(w, r, "auth required", http.StatusUnauthorized)
		return -1
	}

	cred := data.(UserClaims)
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

func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// generateCSRFToken returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateCSRFToken(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func NewCSRFToken(w http.ResponseWriter) error {
	token, err := generateCSRFToken(settings.SecureSettings.CSRFTokenLen)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * settings.SecureSettings.CSRFTokenTTL),
		HttpOnly: false,
		Path:     "/", // TODO: check it
	})

	return nil
}
