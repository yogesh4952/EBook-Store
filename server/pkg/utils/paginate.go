package utils

import (
	"gorm.io/gorm"
)

func Paginate(p Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.Getlimit())
	}
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *Pagination) Getlimit() int {
	switch {
	case p.Limit > 100:
		return 100
	case p.Limit <= 0:
		return 10
	}
	return int(p.Limit)
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.Getlimit()
}
