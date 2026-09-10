package model

type Cart struct {
	Id         uint `json:"id"`
	UserId     uint `json:"user_id"`
	TotalPrice int64

	Items []CartItem
}
