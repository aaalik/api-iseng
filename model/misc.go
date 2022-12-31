package model

type Order int32

const (
	ORDER_UNSPECIFIED Order = iota
	ORDER_ASC
	ORDER_DESC
)

type Range struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type SortBy struct {
	Field string `json:"field,omitempty"`
	Order Order  `json:"order,omitempty"`
}

type Queries struct {
	Page    int32   `json:"page,omitempty"`
	Rows    int32   `json:"rows,omitempty"`
	Sort    *SortBy `json:"sort,omitempty"`
	Keyword string  `json:"keyword,omitempty"`
}
