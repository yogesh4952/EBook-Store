package models

import "time"

type Role string

// 2. Define the allowed enum constants
const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	RoleVendor   Role = "vendor"
)

type User struct {
	Firstname string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	Age       uint      `json:"age"`
	Createdat time.Time `json:"created_at"`
}
