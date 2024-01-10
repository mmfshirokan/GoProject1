package server

import (
	"context"
	"runtime"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/proto/rpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	user repository.RepositoryInterface
	rpc.UnimplementedUserServer
}

func NewUserServer(usr repository.RepositoryInterface) rpc.UserServer {
	return &UserServer{
		user: usr,
	}
}

func (serv *UserServer) GetTroughID(ctx context.Context, req *rpc.RequestGetTroughID) (*rpc.ResponseGetTroughID, error) {
	usr, err := serv.user.GetTroughID(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	return &rpc.ResponseGetTroughID{
		Data: &rpc.UserData{
			Id:   int64(usr.ID),
			Name: usr.Name,
			Male: usr.Male,
		},
	}, nil
}

func (serv *UserServer) Create(ctx context.Context, req *rpc.RequestCreate) (*emptypb.Empty, error) {
	err := serv.user.Create(ctx, model.User{
		ID:   int(req.GetData().Id),
		Name: req.GetData().Name,
		Male: req.GetData().Male,
	})
	if err != nil {
		logError(err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (serv *UserServer) Update(ctx context.Context, req *rpc.RequestUpdate) (*emptypb.Empty, error) {
	err := serv.user.Update(ctx, model.User{
		ID:   int(req.GetData().Id),
		Name: req.GetData().Name,
		Male: req.GetData().Male,
	})
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (serv *UserServer) Delete(ctx context.Context, req *rpc.RequestDelete) (*emptypb.Empty, error) {
	err := serv.user.Delete(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func logError(err error) {
	pc, file, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)

	if ok && details != nil {
		log.WithFields(log.Fields{
			"file":   file,
			"method": details.Name(),
			"line":   line,
			"error":  err,
		}).Error("Error ocured in rpc!")

		return
	}

	log.Fatal("fatal loger error; runtime can't execute Caller")
}
