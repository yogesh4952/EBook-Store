package models

import (
	baseModels "github.com/yogesh4952/ebookstore/pkg/models"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	Roleseller   Role = "seller"
)

type User struct {
	baseModels.Base
	Firstname   string `json:"first_name" `
	Lastname    string `json:"last_name"`
	Email       string `json:"email" gorm:"unique;not null"`
	Role        Role   `json:"role" gorm:"check:role IN ('admin','seller','customer');not null;default:'customer'"`
	PhoneNumber string `json:"phone_number" gorm:"not null; unique"`
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, Roleseller, RoleCustomer:
		return true
	}
	return false
}
