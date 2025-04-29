package dto

type FieldRequest struct {
	Name        string `json:"Name" form:"Name" validate:"required"`
	Type        string `json:"Type" form:"Type" validate:"required,min=3,max=50"`
	Description string `json:"Description" form:"Description" validate:"required,max=50"`
	Price       uint32 `json:"Price" form:"Price" valdiate:"required,gt=0"`
}

type DtoField struct {
	Id          uint   `validate:"required"`
	Name        string `validate:"required,min=3,max=50"`
	Type        string `validate:"required,max=50"`
	Description string `valdiate:"required"`
	Price       uint32 `valdiate:"required,gt=0"`
}
