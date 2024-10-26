package utils

import (
	"gorm.io/gorm"
)

const (
	defaultLimit = 100
	defaultPage  = 1
)

type Pagination struct {
	Limit int
	Page  int
}

func NewPaginate(limit, page int) *Pagination {
	return &Pagination{Limit: limit, Page: page}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	return p.Page
}

func (p *Pagination) PaginatedResult(db *gorm.DB) *gorm.DB {
	offset := p.GetOffset()

	return db.Offset(offset).Limit(p.Limit)
}
