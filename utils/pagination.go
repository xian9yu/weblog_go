package utils

// Pagination 分页
type Pagination struct {
	PageNo   int `form:"page"` //页码
	PageSize int `form:"size"` //每页显示数量
}

func (p Pagination) Offset() int {
	offset := (p.PageNo - 1) * p.PageSize
	if offset < 0 {
		offset = 0
	}
	return offset
}

// 每页默认显示数量 10
func (p Pagination) DefaultPageSize(pageSize int) int {
	if pageSize == 0 {
		pageSize = 10
	}
	return pageSize
}
