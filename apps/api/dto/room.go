package dto

import "github.com/folklinoff/hack_and_change/models"

type CreateRoomRequestDTO struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Password string `json:"password"`
}

type JoinRoomRequestDTO struct {
	Password string `json:"password"`
}

func CreateRoomRequestDTOToModel(dto CreateRoomRequestDTO) models.Room {
	return models.Room{
		Name:     dto.Name,
		Private:  dto.Private,
		Password: dto.Password,
	}
}
