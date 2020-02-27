package utils

import (
	"failless/db"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// Create the JWT key used to create the signature
// todo: rewrite to env variables
var jwtKey = []byte("removeBeforeDebug")

func createJWTToken(user db.User) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(time.Hour * 24 * 30)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
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
	http.SetCookie(w, &http.Cookie{
		Name:     "name",
		Value:    user.Name,
		Expires:  expires,
		HttpOnly: false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "uid",
		Value:    strconv.Itoa(user.Uid),
		Expires:  expires,
		HttpOnly: false,
	})

	return nil
}

func CreateLogout(w http.ResponseWriter) error {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "-",
		MaxAge:   0,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "name",
		Value:    "-",
		MaxAge:   0,
		HttpOnly: false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "uid",
		Value:    "-",
		MaxAge:   0,
		HttpOnly: false,
	})

	return nil
}

func IsAuth(w http.ResponseWriter, r *http.Request) (int, error) {
	c, err := r.Cookie("token")
	if err != nil {
		log.Print(err.Error())
		return -1, nil
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return -1, nil
		}
		w.WriteHeader(http.StatusBadRequest)
		return -1, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return -1, nil
	}
	c, err = r.Cookie("uid")
	if err != nil {
		log.Print(err.Error())
		return -1, nil
	}

	// Get the JWT string from the cookie
	uid, err := strconv.Atoi(c.Value)
	return uid, err
}

func InfoFromCookie(r *http.Request) (info db.User, err error) {
	c, err := r.Cookie("name")
	if err != nil {
		log.Print(err.Error())
		return
	}
	info.Name = c.Value
	c, err = r.Cookie("uid")
	if err != nil {
		log.Print(err.Error())
		return
	}
	info.Uid, err = strconv.Atoi(c.Value)
	if err != nil {
		log.Print(err.Error())
		return
	}
	return info, nil
}