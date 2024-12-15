package domain

//go:generate easyjson -all response.go

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error"`
}

//easyjson:json
type ResponseMessage struct {
	Message string `json:"message"`
}

//easyjson:json
type WrongFieldErrorResponse struct {
	Error       string   `json:"error"`
	WrongFields []string `json:"wrongFields"`
}
