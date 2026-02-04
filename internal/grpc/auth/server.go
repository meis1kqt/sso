package auth

import (
	"context"
	"log/slog"

	v1pb "github.com/meis1kqt/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int64) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64)(bool, error)
}

type serverApi struct {
	v1pb.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	v1pb.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) Login(ctx context.Context, req *v1pb.LoginRequest) (*v1pb.LoginResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == int64(emptyValue) {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		slog.Error("login error", "error", err)
		return nil, err
	}
	return &v1pb.LoginResponse{Token: token}, nil
}

func (s *serverApi) Register(ctx context.Context, req *v1pb.RegisterRequest) (*v1pb.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, err
	}
	return &v1pb.RegisterResponse{UserId: userID}, nil
}

func (s *serverApi) IsAdmin(ctx context.Context, req *v1pb.IdAdminRequest) (*v1pb.IdAdminResponse, error) {
	
	if req.GetUserId() == int64(emptyValue) {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	IsAmdin , err := s.auth.IsAdmin(ctx, req.GetUserId()) 

	if err != nil {
		return nil, err
	}

	return &v1pb.IdAdminResponse{IsAdmin: IsAmdin}, nil
 }
