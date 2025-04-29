package grpcdelivery

import (
	"context"

	userpb "github.com/DevisArya/learn-microservices-protorepo/pb/user"
	"github.com/DevisArya/learn-microservices/user-service/internal/dto"
	"github.com/DevisArya/learn-microservices/user-service/internal/entity"
	"github.com/DevisArya/learn-microservices/user-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController interface {
	userpb.UserServiceServer
}

type UserControllerImpl struct {
	userpb.UnimplementedUserServiceServer
	userUC usecase.UserService
}

func NewUserController(userUc usecase.UserService) UserController {
	return &UserControllerImpl{
		userUC: userUc,
	}
}

func (controller *UserControllerImpl) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {

	userCreateReq := dto.UserCreateRequest{
		Email:        req.GetEmail(),
		Name:         req.GetName(),
		Password:     req.GetPassword(),
		PhoneNumbner: req.GetPassword(),
	}
	id, err := controller.userUC.Create(ctx, &userCreateReq, entity.RoleUser)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userpb.CreateUserResponse{
		Id: &userpb.Id{
			Id: uint32(*id),
		},
	}, nil
}

func (controller *UserControllerImpl) GetUser(ctx context.Context, req *userpb.Id) (*userpb.GetUserResponse, error) {

	user, err := controller.userUC.FindById(ctx, uint(req.GetId()))

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:          req,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}
