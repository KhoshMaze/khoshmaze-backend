package common

type Pagination struct {
	Page     int
	PageSize int
}

type PaginatedResponse[T any] struct {
	Items      []T
	TotalItems int64
	Page       int
	PageSize   int
	TotalPages int
}

func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func NewPaginatedResponse[T any](items []T, totalItems int64, page, pageSize int) *PaginatedResponse[T] {
	totalPages := (int(totalItems) + pageSize - 1) / pageSize
	return &PaginatedResponse[T]{
		Items:      items,
		TotalItems: totalItems,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
