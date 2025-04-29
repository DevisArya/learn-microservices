package dto

type UserCreateRequest struct {
	Email        string `json:"email" form:"email" validate:"required,email,max=255"`
	Name         string `json:"name" form:"name" validate:"required,min=4,max=255"`
	Password     string `json:"password" form:"password" validate:"required,min=8,max=255"`
	PhoneNumbner string `json:"phoneNumber" form:"phoneNumber" validate:"required,min=8,max=20,numeric"`
}

type UserupdatePasswordRequest struct {
	Password string `json:"password" form:"password" validate:"required,min=8,max=255"`
}
type UserupdateEmailRequest struct {
	Email string `json:"email" form:"email" validate:"required,email,max=255"`
}
type UserUpdateProfileRequest struct {
	Name         string `json:"name" form:"name" validate:"required,min=4,max=255"`
	PhoneNumbner string `json:"phoneNumber" form:"phoneNumber" validate:"required,min=8,max=20,numeric"`
}

type OperatorUpdateRequest struct {
	Email        string `json:"email" form:"email" validate:"required,email,max=255"`
	Name         string `json:"name" form:"name" validate:"required,min=4,max=255"`
	Password     string `json:"password" form:"password" validate:"required,min=8,max=255"`
	PhoneNumbner string `json:"phoneNumber" form:"phoneNumber" validate:"required,min=8,max=20,numeric"`
}
