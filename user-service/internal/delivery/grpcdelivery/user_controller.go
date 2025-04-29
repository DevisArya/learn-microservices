package grpcdelivery

import (
	"context"

	"github.com/DevisArya/learn-microservices-protorepo/pb/pagination"
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
	userUC usecase.UserUseCase
}

func NewUserController(userUc usecase.UserUseCase) UserController {
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

func (controller *UserControllerImpl) UpdatePasswordUser(ctx context.Context, req *userpb.UpdatePasswordUserRequest) (*userpb.StatusResponse, error) {

	updatedData := &dto.UserupdatePasswordRequest{
		Password: req.GetPassword(),
	}
	if err := controller.userUC.UpdatePassword(ctx, updatedData, uint(req.Id.GetId())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.StatusResponse{
		Message: "Success update password",
	}, nil
}

func (controller *UserControllerImpl) UpdateEmailUser(ctx context.Context, req *userpb.UpdateEmailUserRequest) (*userpb.StatusResponse, error) {

	updatedData := &dto.UserupdateEmailRequest{
		Email: req.GetEmail(),
	}
	if err := controller.userUC.UpdateEmail(ctx, updatedData, uint(req.Id.GetId())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.StatusResponse{
		Message: "Success update email",
	}, nil
}

func (controller *UserControllerImpl) UpdateProfileUser(ctx context.Context, req *userpb.UpdateProfileUserRequest) (*userpb.StatusResponse, error) {
	updatedData := &dto.UserUpdateProfileRequest{
		Name:         req.GetName(),
		PhoneNumbner: req.GetPhoneNumber(),
	}

	if err := controller.userUC.UpdateProfile(ctx, updatedData, uint(req.Id.GetId())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.StatusResponse{
		Message: "Success update profile",
	}, nil
}

func (controller *UserControllerImpl) Delete(ctx context.Context, req *userpb.Id) (*userpb.StatusResponse, error) {

	if err := controller.userUC.Delete(ctx, uint(req.GetId())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.StatusResponse{
		Message: "Success delete user",
	}, nil
}

// FindAll implements UserHandler
func (controller *UserControllerImpl) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {

	res, paging, err := controller.userUC.FindAll(ctx, req.GetLimit(), req.GetPage())

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var users []*userpb.User
	for _, val := range *res {
		users = append(users, &userpb.User{
			Id:          &userpb.Id{Id: uint32(val.Id)},
			Name:        val.Name,
			Email:       val.Email,
			Password:    val.Email,
			PhoneNumber: val.PhoneNumber,
		})
	}

	return &userpb.GetUsersResponse{
		Pagination: &pagination.Pagination{
			CurrentPage: paging.CurrentPage,
			Limit:       paging.Limit,
			TotalRecord: paging.TotalRecord,
			TotalPage:   paging.TotalPage,
		}, Data: users,
	}, nil
}
