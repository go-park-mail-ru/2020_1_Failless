package security

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"github.com/dgrijalva/jwt-go"
	"log"
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
var jwtKey = []byte("removeBeforeDebug")

func createJWTToken(user db.User) (string, error) {
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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateAuth(w http.ResponseWriter, user db.User) error {
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

func IsAuth(w http.ResponseWriter, r *http.Request) (forms.SignForm, error) {
	c, err := r.Cookie("token")
	if err != nil {
		log.Print(err.Error())
		return forms.SignForm{}, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	log.Println(claims.Name, claims.Phone, claims.Email)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return forms.SignForm{}, err
		}
		//w.WriteHeader(network.StatusBadRequest)
		return forms.SignForm{}, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return forms.SignForm{}, err
	}
	return forms.SignForm{
		Uid: claims.Uid,
		Phone: claims.Phone,
		Email: claims.Email,
		Name: claims.Name}, nil
}
