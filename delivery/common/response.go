package common

//DefaultResponse default payload response
type DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponseWithData(data interface{}) ResponseSuccess {
	return ResponseSuccess{
		Code:    200,
		Message: "Successful Operation",
		Data:    data,
	}
}

func ErrorResponse(code int, message string) DefaultResponse {
	return DefaultResponse{
		Code:    code,
		Message: message,
	}
}

//NewSuccessOperationResponse default success operation response
func NewSuccessOperationResponse() DefaultResponse {
	return DefaultResponse{
		200,
		"Successful Operation",
	}
}

//NewNotFoundResponse default not found error response
func NewNotFoundResponse() DefaultResponse {
	return DefaultResponse{
		404,
		"Not Found",
	}
}

//NewBadRequestResponse default bad request error response
func NewBadRequestResponse() DefaultResponse {
	return DefaultResponse{
		400,
		"Bad Request",
	}
}
