package models

import "time"

type User struct {
	Firstname string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Role      string    `json:"role"`
	Age       uint      `json:"age"`
	Createdat time.Time `json:"created_at"`
}
