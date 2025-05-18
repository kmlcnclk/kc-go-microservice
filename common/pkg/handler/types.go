package handler

type Request any

type Response any

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}
