package dto

type CreateRoomRequestDTO struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Password string `json:"password"`
}

type JoinRoomRequestDTO struct {
	Password string `json:"password"`
}
