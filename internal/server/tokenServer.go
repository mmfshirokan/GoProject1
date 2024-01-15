package server

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/internal/service"
	"github.com/mmfshirokan/GoProject1/proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TokenServer struct {
	user     repository.RepositoryInterface
	password repository.PwRepositoryInterface
	token    repository.AuthRepositoryInterface
	pb.UnimplementedTokenServer
}

func NewTokenServer(usr repository.RepositoryInterface, pwd repository.PwRepositoryInterface, tok repository.AuthRepositoryInterface) pb.TokenServer {
	return &TokenServer{
		user:     usr,
		password: pwd,
		token:    tok,
	}
}

func (serv *TokenServer) SignUp(ctx context.Context, req *pb.RequestSignUp) (*emptypb.Empty, error) {
	err := serv.user.Create(ctx, model.User{
		ID:   int(req.GetData().GetId()),
		Name: req.GetData().GetName(),
		Male: req.GetData().GetMale(),
	})
	if err != nil {
		logError(err)
		return nil, err
	}

	err = serv.password.Store(ctx, model.User{
		ID:       int(req.GetData().GetId()),
		Password: req.GetPassword(),
	})
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (serv *TokenServer) SignIn(ctx context.Context, req *pb.RequestSignIn) (*pb.ResponseSignIn, error) {
	valid, err := serv.password.Compare(ctx, model.User{
		ID:       int(req.GetUserID()),
		Password: req.GetPassword(),
	})
	if err != nil {
		logError(err)
		return nil, err
	}
	if !valid {
		err = errors.New("wrong password")
		logError(err)
		return nil, err
	}

	user, err := serv.user.GetTroughID(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	authToken := service.CreateAuthToken(
		user.ID,
		user.Name,
		user.Male,
	)

	refreshToken, err := newRFT(user.ID)
	if err != nil {
		logError(err)
		return nil, err
	}

	err = serv.token.Create(ctx, refreshToken)
	if err != nil {
		logError(err)
		return nil, err
	}

	return &pb.ResponseSignIn{
		Tokens: &pb.Jwt{
			AuthToken: authToken,
			Rft: &pb.RefreshToken{
				UserID:     int64(refreshToken.UserID),
				Uuid:       refreshToken.ID.String(),
				Hash:       refreshToken.Hash,
				Experation: timestamppb.New(refreshToken.Expiration),
			},
		},
	}, nil
}

func (serv *TokenServer) Refresh(ctx context.Context, req *pb.RequestRefresh) (*pb.ResponseRefresh, error) {
	id, err := uuid.Parse(req.GetRft().GetUuid())
	userID := int(req.GetRft().GetUserID())
	if err != nil {
		logError(err)
		return nil, err
	}

	valid, err := service.ValidateRfTokenTroughID(req.GetRft().GetHash(), id)
	if err != nil {
		logError(err)
		return nil, err
	}

	if !valid {
		err = errors.New("invalid rft")
		logError(err)
		return nil, err
	}

	err = serv.token.Delete(ctx, id)
	if err != nil {
		logError(err)
		return nil, err
	}

	user, err := serv.user.GetTroughID(ctx, userID)
	if err != nil {
		logError(err)
		return nil, err
	}

	authToken := service.CreateAuthToken(userID, user.Name, user.Male)

	refreshToken, err := newRFT(userID)
	if err != nil {
		logError(err)
		return nil, err
	}

	err = serv.token.Create(ctx, refreshToken)
	if err != nil {
		logError(err)
		return nil, err
	}

	return &pb.ResponseRefresh{
		Tokens: &pb.Jwt{
			AuthToken: authToken,
			Rft: &pb.RefreshToken{
				UserID:     int64(refreshToken.UserID),
				Uuid:       refreshToken.ID.String(),
				Hash:       refreshToken.Hash,
				Experation: timestamppb.New(refreshToken.Expiration),
			},
		},
	}, nil
}

// type TokenServer interface {
//     SignUp(context.Context, *RequestSignUp) (*emptypb.Empty, error)
//     SignIn(context.Context, *RequestSignIn) (*ResponseSignIn, error)
//     Refresh(context.Context, *RequestRefresh) (*ResponseRefresh, error)
//     mustEmbedUnimplementedTokenServer()
// }

// suplemental function
func newRFT(userID int) (*model.RefreshToken, error) {
	const refreshTokenLifeTime = 12

	id := uuid.New()
	hashedID, err := service.ConductHashing(id)
	if err != nil {
		return nil, err
	}

	return &model.RefreshToken{
		UserID:     userID,
		ID:         id,
		Hash:       hashedID,
		Expiration: time.Now().Add(time.Hour * refreshTokenLifeTime),
	}, nil
}
