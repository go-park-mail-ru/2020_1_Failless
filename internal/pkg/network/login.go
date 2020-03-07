package network

import (
	"failless/internal/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Uid   int    `json:"uid"`
	jwt.StandardClaims
}

// Create the JWT key used to create the signature
// todo: rewrite to env variables
var JwtKey = []byte("removeAfterDebug")

func createJWTToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * 30)
	claims := &Claims{
		Email: user.Email,
		Phone: user.Phone,
		Name:  user.Name,
		Uid:   user.Uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateAuth(w http.ResponseWriter, user models.User) error {
	token, err := createJWTToken(user)
	if err != nil {
		return err
	}

	expires := time.Now().Add(time.Hour * 24 * 30)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expires,
		HttpOnly: true,
	})
	return nil
}

func CreateLogout(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "-",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return nil
}
