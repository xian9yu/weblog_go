package utils

// Pagination 分页
type Pagination struct {
	PageNo    int `form:"page"`
	LimitSize int `form:"size"`
}

func (p Pagination) Offset() int {
	offset := (p.PageNo - 1) * p.LimitSize
	if offset < 0 {
		offset = 0
	}
	return offset
}
