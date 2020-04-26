package usecase

import (
	"context"
	"failless/internal/pkg/auth"
	"failless/internal/pkg/auth/repository"
	"failless/internal/pkg/db"
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

func GetUseCase() AuthService {
	return AuthService{
		Rep: repository.NewSqlAuthRepository(db.ConnectToDB()),
	}
}

func (as *AuthService) Authorize(ctx context.Context, cred *pb.AuthRequest) (*pb.AuthorizeReply, error) {
	log.Print("Authorize: ")
	user, err := as.Rep.GetUserByPhoneOrEmail(cred.Phone, cred.Email)
	if err == nil && user.Uid < 0 {
		log.Println("client - user not found")
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: "User doesn't exist",
		}, nil
	} else if err != nil {
		log.Println("error - ", err.Error())
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: err.Error(),
		}, err
	}

	if !security.ComparePasswords(user.Password, cred.Password) {
		log.Println("client - passwords are not equal")
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: "Passwords are not equal",
		}, nil
	}

	token, err := network.CreateJWTToken(user)
	if err != nil {
		log.Println("error - ", err.Error())
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: err.Error(),
		}, nil
	}
	log.Println("OK")
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
	log.Print("CheckAuthorize: ")
	// Get the JWT string from the cookie
	tknStr := in.Token
	claims := &network.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return network.JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("client - token signature invalid")
			return &pb.AuthorizeReply{
				Ok:      false,
				Cred:    nil,
				Message: "invalid token",
			}, nil
		}
		log.Println("error - ", err.Error())
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: err.Error(),
		}, err
	}

	if !tkn.Valid {
		log.Println("client - token invalid")
		return &pb.AuthorizeReply{
			Ok:      false,
			Cred:    nil,
			Message: "invalid token",
		}, nil
	}
	// success. user is authorized
	log.Println("OK")
	ctx = context.WithValue(ctx, security.CtxUserKey, claims)

	return &pb.AuthorizeReply{
		Ok: true,
		Cred: &pb.Credentials{
			Uid:   int64(claims.Uid),
			Phone: claims.Phone,
			Email: claims.Email,
			Name:  claims.Name,
		},
		Message: "",
	}, nil

}

func (*AuthService) GetToken(ctx context.Context, in *pb.AuthRequest) (*pb.Token, error) {
	log.Print("GetToken: ")
	user := models.User{
		Uid:   int(in.Uid),
		Name:  in.Name,
		Phone: in.Phone,
		Email: in.Email,
	}
	pass, err := security.EncryptPassword(in.Password)
	if err != nil {
		log.Println("error - ", err.Error())
		return nil, err
	}

	user.Password = pass
	token, err := network.CreateJWTToken(user)
	if err != nil {
		log.Println("error - ", err.Error())
		return nil, err
	}
	log.Println("OK")
	return &pb.Token{Token: token}, err
}
