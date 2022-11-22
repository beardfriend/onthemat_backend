package utils

type Pagination struct {
	offset   int
	limit    int
	total    int
	pageSize int
	pageNo   int
}

func NewPagination(pageNo, pageSize int) *Pagination {
	return &Pagination{
		pageSize: pageSize,
		pageNo:   pageNo,
		limit:    pageSize,
		offset:   (pageNo - 1) * pageSize,
	}
}

func (p *Pagination) SetTotal(total int) {
	p.total = total
}

func (p *Pagination) GetLimit() int {
	return p.limit
}

func (p *Pagination) GetOffset() int {
	return p.offset
}

type PagenationInfo struct {
	PageSize  int
	PageNo    int
	PageCount int
	RowCount  int
}

func (p *Pagination) GetInfo(resultLength int) *PagenationInfo {
	pageCount := p.total/p.pageSize + 1
	if p.total%p.pageSize == 0 {
		pageCount = p.total / p.pageSize
	}

	return &PagenationInfo{
		PageSize:  p.pageSize,
		PageNo:    p.pageNo,
		PageCount: pageCount,
		RowCount:  resultLength,
	}
}
