package dto

import "github.com/folklinoff/hack_and_change/models"

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequestDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserModelResponseDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func LoginRequestDTOToUserModel(req LoginRequestDTO) models.User {
	return models.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func SignupRequestDTOToUserModel(req SignupRequestDTO) models.User {
	return models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Password,
	}
}

func UserModelToUserModelResponseDTO(req models.User) UserModelResponseDTO {
	return UserModelResponseDTO{
		ID:    req.ID,
		Email: req.Email,
		Name:  req.Name,
	}
}
