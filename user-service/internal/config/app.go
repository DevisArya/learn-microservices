package config

import (
	"net"

	userpb "github.com/DevisArya/learn-microservices-protorepo/pb/user"
	"github.com/DevisArya/learn-microservices/user-service/internal/delivery/grpcdelivery"
	"github.com/DevisArya/learn-microservices/user-service/internal/repository"
	"github.com/DevisArya/learn-microservices/user-service/internal/usecase"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

type BootstrapResult struct {
	GRPCServer *grpc.Server
	Listener   net.Listener
}

func Bootstrap(cfg *BootstrapConfig) (*BootstrapResult, error) {

	// Setup TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return nil, err
	}

	//Init depedencies
	fieldRepo := repository.NewUserRepository()
	fieldUc := usecase.NewUserUseCase(fieldRepo, cfg.DB, cfg.Validate)
	fieldCtrl := grpcdelivery.NewUserController(fieldUc)

	//init grpc server & register service

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, fieldCtrl)

	return &BootstrapResult{
		GRPCServer: grpcServer,
		Listener:   lis,
	}, nil
}
