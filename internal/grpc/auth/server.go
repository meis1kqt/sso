package auth

import (
	"context"

	v1pb "github.com/meis1kqt/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverApi struct {
	v1pb.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	v1pb.RegisterAuthServer(gRPC, &serverApi{})
}

func (s *serverApi) Login(ctx context.Context, req *v1pb.LoginRequest) (*v1pb.LoginResponse, error) {
	
	return &v1pb.LoginResponse{Token: "lashfsafshjfsdkfdjkadsfkPIDORS"}, nil
}

func (s *serverApi) Register(ctx context.Context, req *v1pb.RegisterRequest) (*v1pb.RegisterResponse, error) {
	panic("reg test")
}

func (s *serverApi) IsAdmin(ctx context.Context, req *v1pb.IdAdminRequest) (*v1pb.IdAdminResponse, error) {
	panic("reg test")
}
