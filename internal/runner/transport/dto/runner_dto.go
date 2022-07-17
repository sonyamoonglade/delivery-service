package dto

type RegisterRunnerDto struct {
	Username    string `json:"username" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
