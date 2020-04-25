package middleware

import (
	"context"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"failless/internal/pkg/settings"
	"log"
	"net/http"

	pb "failless/api/proto/auth"
)

// Structure for describe an error
type authError struct {

	// Error message
	msg string

	// 1 - cookie not found
	// 2 - parse error
	// 3 - signature invalid
	// 4 - token invalid
	code int
}

//type UserClaims struct {
//	Uid   int
//	Phone string
//	Email string
//	Name  string
//}

//type UserKey string
//
//// Context variable for pushing credentials through middleware to handlers
//const CtxUserKey UserKey = "auth"

// Auth middleware checks is user authorized
// If user is not authorized it write failed checker to the authError structure
func Auth(next settings.HandlerFunc) settings.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		ctx := context.Background()
		c, err := r.Cookie("token")

		// error - no cookie was found
		if err != nil {
			ctx = context.WithValue(ctx, security.CtxUserKey, nil)
			log.Print(err.Error())
			network.GenErrorCode(w, r, err.Error(), http.StatusUnauthorized)
			return
		} else {
			authReply, err := settings.AuthClient.CheckAuthorize(
				ctx, &pb.Token{Token: c.Value})
			if err != nil {
				ctx = context.WithValue(ctx, security.CtxUserKey, nil)
				w.WriteHeader(http.StatusInternalServerError)
				network.GenErrorCode(w, r, err.Error(), http.StatusInternalServerError)
				return
			}

			uClaims := security.UserClaims{
				Uid:   int(authReply.Cred.Uid),
				Phone: authReply.Cred.Phone,
				Email: authReply.Cred.Email,
				Name:  authReply.Cred.Name,
			}
			ctx = context.WithValue(ctx, security.CtxUserKey, uClaims)

		}
		next(w, r.WithContext(ctx), ps)
	}
}
