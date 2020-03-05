package security

import (
	"failless/internal/pkg/db"
	"failless/internal/pkg/forms"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)



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
		//w.WriteHeader(middleware.StatusBadRequest)
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
