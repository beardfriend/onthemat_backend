package utils

type pagination struct {
	offset int
	limit  int
	total  int
}

func NewPagination(pageNo, pageSize int) *pagination {
	return &pagination{
		limit:  pageSize,
		offset: (pageNo - 1) * pageSize,
	}
}

func (p *pagination) SetTotal(total int) {
	p.total = total
}

func (p *pagination) GetLimit() int {
	return p.limit
}

func (p *pagination) GetOffset() int {
	return p.offset
}
