package dto

type PageQuery struct {
	Page     int
	PageSize int
}

type PageResponse[T any] struct {
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	Data     []*T `json:"data"`
}
