package models

type Room struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	OwnerId  int64     `json:"owner_id"`
	Private  bool      `json:"private"`
	Password string    `json:"password"`
	Users    []User    `json:"user_ids"`
	Messages []Message `json:"messages"`
}
