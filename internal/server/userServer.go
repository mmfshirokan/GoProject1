package server

import (
	"context"
	"runtime"

	"github.com/mmfshirokan/GoProject1/internal/model"
	"github.com/mmfshirokan/GoProject1/internal/repository"
	"github.com/mmfshirokan/GoProject1/proto/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	user repository.RepositoryInterface
	pb.UnimplementedUserServer
}

func NewUserServer(usr repository.RepositoryInterface) pb.UserServer {
	return &UserServer{
		user: usr,
	}
}

func (serv *UserServer) GetUser(ctx context.Context, req *pb.RequestGetUser) (*pb.ResponseGetUser, error) {
	ctx = newMetadataContext(ctx, req.AuthToken)

	user, err := serv.user.GetTroughID(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	return &pb.ResponseGetUser{
		Data: &pb.UserData{
			Id:   int64(user.ID),
			Name: user.Name,
			Male: user.Male,
		},
	}, nil
}

func (serv *UserServer) UpdateUser(ctx context.Context, req *pb.RequestUpdateUser) (*emptypb.Empty, error) { // NOTE mabe add option to update password
	ctx = newMetadataContext(ctx, req.AuthToken)

	err := serv.user.Update(ctx, model.User{
		//ID:   int(req.GetData().GetId()),
		Name: req.GetData().GetName(),
		Male: req.GetData().GetMale(),
	})
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (serv *UserServer) DeleteUser(ctx context.Context, req *pb.RequestDelete) (*emptypb.Empty, error) {
	ctx = newMetadataContext(ctx, req.GetAuthToken())

	err := serv.user.Delete(ctx, int(req.GetUserID()))
	if err != nil {
		logError(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// suplimental functions method

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

func newMetadataContext(ctx context.Context, auth string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", auth))
}
