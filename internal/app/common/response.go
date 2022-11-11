package common

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseWithData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type ResponseWithPagination struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
	Pagination interface{} `json:"pagination"`
}
