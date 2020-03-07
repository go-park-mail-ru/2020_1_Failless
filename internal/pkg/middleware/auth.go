package middleware

import (
	"context"
	"failless/internal/pkg/forms"
	"failless/internal/pkg/network"
	"failless/internal/pkg/settings"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

// Structure for describe an error
type authError struct {

	// Error message
	msg  string

	// 1 - cookie not found
	// 2 - parse error
	// 3 - signature invalid
	// 4 - token invalid
	code int
}

type UserKey string

// Context variable for pushing credentials through middleware to handlers
const CtxUserKey UserKey = "auth"

// Auth middleware checks is user authorized
// If user is not authorized it write failed checker to the authError structure
func Auth(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		var errMsg authError
		ctx := context.Background()
		c, err := r.Cookie("token")

		// error - no cookie was found
		if err != nil {
			ctx = context.WithValue(ctx, CtxUserKey, nil)
			errMsg.code = 1
			errMsg.msg = err.Error()
			log.Print(err.Error())
		} else {

			// Get the JWT string from the cookie
			tknStr := c.Value
			claims := &network.Claims{}
			tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
				return network.JwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					errMsg.code = 3
					ctx = context.WithValue(ctx, CtxUserKey, nil)
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					errMsg.code = 5
				}
				errMsg.msg = err.Error()
			}

			if errMsg.code == 0 {
				if !tkn.Valid {
					ctx = context.WithValue(ctx, "auth", nil)
					w.WriteHeader(http.StatusUnauthorized)
					errMsg.code = 5
					errMsg.msg = "Token invalid"
				} else { // success. user is authorized
					form := forms.SignForm{
						Uid:   claims.Uid,
						Phone: claims.Phone,
						Email: claims.Email,
						Name:  claims.Name,
					}
					ctx = context.WithValue(ctx, CtxUserKey, form)
				}
			}
		}
		next(w, r.WithContext(ctx), ps)
	}
}
