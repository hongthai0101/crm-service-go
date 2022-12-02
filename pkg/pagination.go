package pkg

type Pagination[T interface{}] struct {
	Limit int64 `json:"limit"`
	Skip  int64 `json:"skip"`
	Total int64 `json:"total"`
	List  []*T  `json:"list"`
}
