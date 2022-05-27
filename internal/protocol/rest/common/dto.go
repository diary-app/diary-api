package common

type Pagination struct {
	offset int
	count  int
	limit  int
}

type ErrorResponse struct {
	Message string `json:"message"`
}
