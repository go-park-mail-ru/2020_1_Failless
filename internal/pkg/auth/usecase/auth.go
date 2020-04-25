package usecase

import (
	"context"
	"failless/internal/pkg/auth"
	"failless/internal/pkg/models"
	"failless/internal/pkg/network"
	"failless/internal/pkg/security"
	"github.com/dgrijalva/jwt-go"
	"log"

	pb "failless/api/proto/auth"
)

type AuthService struct {
	Rep auth.Repository
}

func (as *AuthService) Authorize(ctx context.Context, cred *pb.AuthRequest) (*pb.AuthorizeReply, error) {
	user, err := as.Rep.GetUserByPhoneOrEmail(cred.Phone, cred.Email)
	if err == nil && user.Uid < 0 {
		log.Println("user not found")
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: "User doesn't exist",
		}, nil
	} else if err != nil {
		log.Println("error was occurred")
		log.Println(err.Error())
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: err.Error(),
		}, err
	}

	if !security.ComparePasswords(user.Password, cred.Password) {
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: "Passwords is not equal",
		}, nil
	}

	token, err := network.CreateJWTToken(user)
	if err != nil {
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: err.Error(),
		}, nil
	}
	return &pb.AuthorizeReply{
		Ok: true,
		Cred: &pb.Credentials{
			Uid:   int64(user.Uid),
			Name:  user.Name,
			Phone: user.Phone,
			Email: user.Email,
		},
		Message: token,
	}, nil
}

func (*AuthService) CheckAuthorize(ctx context.Context, in *pb.Token) (*pb.AuthorizeReply, error) {
	// Get the JWT string from the cookie
	tknStr := in.Token
	claims := &network.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return network.JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &pb.AuthorizeReply{
				Ok:      false,
				Cred:    nil,
				Message: "invalid token",
			}, nil
		}
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: err.Error(),
		}, err
	}

	form := pb.Credentials{}
	if !tkn.Valid {
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: "invalid token",
		}, nil
	} else { // success. user is authorized
		form = pb.Credentials{
			Uid:   int64(claims.Uid),
			Phone: claims.Phone,
			Email: claims.Email,
			Name:  claims.Name,
		}
		ctx = context.WithValue(ctx, security.CtxUserKey, form)
	}

	return &pb.AuthorizeReply{
		Ok:      true,
		Cred:    &form,
		Message: "",
	}, nil

}

func (*AuthService) GetToken(ctx context.Context, in *pb.AuthRequest) (*pb.Token, error) {
	user := models.User{
		Uid:   int(in.Uid),
		Name:  in.Name,
		Phone: in.Phone,
		Email: in.Email,
	}
	pass, err := security.EncryptPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user.Password = pass
	token, err := network.CreateJWTToken(user)
	if err != nil {
		return nil, err
	}

	return &pb.Token{Token: token}, err
}
