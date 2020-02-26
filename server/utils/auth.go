package utils

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2020_1_Failless/db"
	"github.com/go-park-mail-ru/2020_1_Failless/server/forms"
	"golang.org/x/crypto/bcrypt"
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

func CreateAuth(w http.ResponseWriter, r *http.Request, user db.User) error {
	token, err := createJWTToken(user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
	})

	return nil
}

func IsAuth(w http.ResponseWriter, r *http.Request) (bool, error) {
	c, err := r.Cookie("token")
	if err != nil {
		log.Print(err.Error())
		//if err == http.ErrNoCookie {
		//	return false, nil
		//}
		return false, nil
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
			return false, nil
		}
		w.WriteHeader(http.StatusBadRequest)
		return false, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	page := claims
	_ = json.NewEncoder(w).Encode(page)
	return true, nil
}

func EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func RegisterNewUser(user forms.SignForm) error {
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
