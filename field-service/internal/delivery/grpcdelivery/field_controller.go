package grpcdelivery

import (
	"context"

	fieldpb "github.com/DevisArya/learn-microservices-protorepo/pb/field"
	pagingpb "github.com/DevisArya/learn-microservices-protorepo/pb/pagination"
	"github.com/DevisArya/learn-microservices/field-service/internal/dto"
	"github.com/DevisArya/learn-microservices/field-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FieldController interface {
	fieldpb.FieldServiceServer
}

type FieldControllerImpl struct {
	FieldUc usecase.FieldUseCase
	fieldpb.UnimplementedFieldServiceServer
}

func NewFieldController(fieldUc usecase.FieldUseCase) FieldController {
	return &FieldControllerImpl{
		FieldUc: fieldUc,
	}
}

func (controller *FieldControllerImpl) GetFields(ctx context.Context, req *fieldpb.GetFieldsRequest) (*fieldpb.GetFieldsResponse, error) {

	res, paging, err := controller.FieldUc.FindAll(ctx, req.GetLimit(), req.GetPage())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var fields []*fieldpb.Field
	for _, f := range *res {
		fields = append(fields, &fieldpb.Field{
			Id:          uint32(f.Id),
			Name:        f.Name,
			Type:        f.Type,
			Price:       uint64(f.Price),
			Description: f.Description,
		})
	}

	return &fieldpb.GetFieldsResponse{
		Pagination: &pagingpb.Pagination{
			CurrentPage: paging.CurrentPage,
			Limit:       paging.Limit,
			TotalRecord: paging.TotalRecord,
			TotalPage:   paging.TotalPage,
		}, Data: fields,
	}, nil
}

func (controller *FieldControllerImpl) GetField(ctx context.Context, req *fieldpb.Id) (*fieldpb.GetFieldResponse, error) {

	res, err := controller.FieldUc.FindById(ctx, uint(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fieldpb.GetFieldResponse{
		Field: &fieldpb.Field{
			Id:          uint32(res.Id),
			Name:        res.Name,
			Type:        res.Type,
			Description: res.Description,
			Price:       uint64(res.Price),
		},
	}, nil
}

func (controller *FieldControllerImpl) CreateField(ctx context.Context, req *fieldpb.CreateFieldRequest) (*fieldpb.CreateFieldResponse, error) {

	fieldReq := dto.FieldRequest{
		Name:        req.GetName(),
		Type:        req.GetType(),
		Description: req.GetDescription(),
		Price:       uint32(req.GetPrice()),
	}
	field, err := controller.FieldUc.Save(ctx, &fieldReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fieldpb.CreateFieldResponse{
		Id:      uint32(field.Id),
		Message: "Create field succesfully",
	}, nil
}

func (controller *FieldControllerImpl) UpdateField(ctx context.Context, req *fieldpb.UpdateFieldRequest) (*fieldpb.StatusResponse, error) {
	fieldReq := dto.FieldRequest{
		Name:        req.GetName(),
		Type:        req.GetType(),
		Description: req.GetDescription(),
		Price:       uint32(req.GetPrice()),
	}

	if err := controller.FieldUc.Update(ctx, &fieldReq, uint(req.GetId())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fieldpb.StatusResponse{
		Message: "Success update",
	}, nil
}
func (controller *FieldControllerImpl) DeleteField(ctx context.Context, req *fieldpb.Id) (*fieldpb.StatusResponse, error) {

	if err := controller.FieldUc.Delete(ctx, uint(req.GetId())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fieldpb.StatusResponse{
		Message: "Succes delete",
	}, nil
}
