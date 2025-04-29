package config

import (
	"net"

	fieldpb "github.com/DevisArya/learn-microservices-protorepo/pb/field"
	"github.com/DevisArya/learn-microservices/field-service/internal/delivery/grpcdelivery"
	"github.com/DevisArya/learn-microservices/field-service/internal/repository"
	"github.com/DevisArya/learn-microservices/field-service/internal/usecase"
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
	fieldRepo := repository.NewFieldRepository()
	fieldUc := usecase.NewFieldUseCase(fieldRepo, cfg.DB, cfg.Validate)
	fieldCtrl := grpcdelivery.NewFieldController(fieldUc)

	//init grpc server & register service

	grpcServer := grpc.NewServer()
	fieldpb.RegisterFieldServiceServer(grpcServer, fieldCtrl)

	return &BootstrapResult{
		GRPCServer: grpcServer,
		Listener:   lis,
	}, nil
}
