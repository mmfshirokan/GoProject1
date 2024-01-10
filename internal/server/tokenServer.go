package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/internal/service"
	"github.com/mmfshirokan/GoProject1/proto/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TokenServer struct {
	token repository.AuthRepositoryInterface
	rpc.UnimplementedTokenServer
}

func NewTokenServer(tok repository.AuthRepositoryInterface) rpc.TokenServer {
	return &TokenServer{
		token: tok,
	}
}

func (serv *TokenServer) CreateRfToken(ctx context.Context, req *rpc.RequestCreateRfToken) (*emptypb.Empty, error) {
	const refreshTokenLifeTime = 12

	id := uuid.New()

	hash, err := service.ConductHashing(id)
	if err != nil {
		logError(err)
		return nil, err
	}

	err = serv.token.Create(ctx, &model.RefreshToken{
		UserID:     int(req.GetUserID()),
		ID:         id,
		Hash:       hash,
		Expiration: time.Now().Add(time.Hour * refreshTokenLifeTime),
	})
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (serv *TokenServer) DeleteRftToken(ctx context.Context, req *rpc.RequestDeleteRftToken) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetUuid())
	if err != nil {
		logError(err)
		return nil, err
	}

	serv.token.Delete(ctx, id)
	return &emptypb.Empty{}, nil
}

func (serv *TokenServer) GetByUserID(ctx context.Context, req *rpc.RequestGetByUserID) (*rpc.ResponseGetByUserID, error) { // TODO add redis
	rfTokens, err := serv.token.GetByUserID(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	var result []*rpc.ResponseGetByUserIDRefreshToken
	for i := range rfTokens {
		result[i] = &rpc.ResponseGetByUserIDRefreshToken{
			UserId:     int64(rfTokens[i].UserID),
			Uuid:       rfTokens[i].ID.String(),
			Hash:       rfTokens[i].Hash,
			Expiration: timestamppb.New(rfTokens[i].Expiration),
		}
	}

	return &rpc.ResponseGetByUserID{
		Rft: result,
	}, nil
}
