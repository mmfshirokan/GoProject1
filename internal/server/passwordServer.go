package server

import (
	"context"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/proto/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PasswordServer struct {
	password repository.PwRepositoryInterface
	rpc.UnimplementedPasswordServer
}

func NewPasswordServer(pwd repository.PwRepositoryInterface) rpc.PasswordServer {
	return &PasswordServer{
		password: pwd,
	}
}

func (serv *PasswordServer) Store(ctx context.Context, req *rpc.RequestStore) (*emptypb.Empty, error) {
	err := serv.password.Store(ctx, model.User{
		ID:       int(req.GetData().UserID),
		Password: req.GetData().Password,
	})
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (serv *PasswordServer) Compare(ctx context.Context, req *rpc.RequestCompare) (*rpc.ResponseCompare, error) {
	res, err := serv.password.Compare(ctx, model.User{
		ID:       int(req.GetData().UserID),
		Password: req.GetData().Password,
	})
	if err != nil {
		logError(err)
		return &rpc.ResponseCompare{
			PwdValid: res,
		}, err
	}

	return &rpc.ResponseCompare{
		PwdValid: res,
	}, nil
}

func (serv *PasswordServer) DeletePassword(ctx context.Context, req *rpc.RequestDeletePassword) (*emptypb.Empty, error) {
	err := serv.password.DeletePassword(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
