package adaptersinboundhttpdto

type CommonErrorResponse struct {
	Error string `json:"error"`
}

type CommonPaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}
