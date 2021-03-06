package security

import (
	"crypto/rand"
	"failless/internal/pkg/network"
	"failless/internal/pkg/settings"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

//easyjson:json
type UserClaims struct {
	Uid   int    `json:"uid"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserKey string

// Context variable for pushing credentials through middleware to handlers
const CtxUserKey UserKey = "auth"

// Compare password from database in bytes format and password
// which was gotten from cookie
func ComparePasswords(hash []byte, p string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(p))
	return err == nil
}

// Generate bcrypt hash with default cost from input string
// Returns bytes array
// TODO: check is default cost ok
func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckCredentials gets from context structure with credentials
// (see UserClaims above) and check is Uid from this structure valid
// if it is valid, than function returns it, else function send to
// the client message with description of error and status code
func CheckCredentials(w http.ResponseWriter, r *http.Request) int {
	data := r.Context().Value(CtxUserKey)
	if data == nil {
		network.GenErrorCode(w, r, network.MessageErrorAuthRequired, http.StatusUnauthorized)
		return -1
	}

	cred := data.(UserClaims)
	if cred.Uid < 0 {
		network.GenErrorCode(w, r, network.MessageErrorIncorrectTokenUid, http.StatusBadRequest)
		return -1
	}

	return cred.Uid
}

func GetUserFromCtx(r *http.Request) (UserClaims, error) {
	data := r.Context().Value(CtxUserKey)
	if data == nil {
		return UserClaims{}, claimsNotFoundError
	}

	cred := data.(UserClaims)
	if cred.Uid < 0 {
		return UserClaims{}, incorrectTokenUidError
	}

	return cred, nil
}

func CompareUidsFromURLAndToken(w http.ResponseWriter, r *http.Request, ps map[string]string) int {
	ctxUid := CheckCredentials(w, r)

	uid := int64(0)
	if uid = network.GetIdFromRequest(w, r, ps); uid < 0 {
		return -1
	}

	if ctxUid != int(uid) {
		network.GenErrorCode(w, r, "forbidden", http.StatusForbidden)
		return -1
	}

	return int(uid)
}

// Simple function that generate random sequence of bytes for
// different aims such as security, fun, fake data
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

// NewCSRFToken set to http.ResponseWriter cookie "csrf" which contains
// just generated csrf-token using generateCSRFToken function. NewCSRFToken
// use constant CSRFTokenLen from settings structure SecureSettings which
// are set in the configs/server/settings.go
func NewCSRFToken(w http.ResponseWriter) error {
	token, err := generateCSRFToken(settings.SecureSettings.CSRFTokenLen)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * settings.SecureSettings.CSRFTokenTTL),
		HttpOnly: false,
		Path:     "/", // TODO: check it
	})

	return nil
}
