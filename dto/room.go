package dto

type CreateRoomDTO struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Password string `json:"password"`
}
