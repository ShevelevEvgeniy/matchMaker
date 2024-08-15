package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Status string `json:"status"`
}

const (
	StatusOK         = "200"
	StatusError      = "500"
	StatusBadRequest = "400"
	StatusConflict   = "409"
)

func OK() SuccessResponse {
	return SuccessResponse{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func BadRequest(msg string) Response {
	return Response{
		Status: StatusBadRequest,
		Error:  msg,
	}
}

func Conflict(msg string) Response {
	return Response{
		Status: StatusConflict,
		Error:  msg,
	}
}

func NotFound(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func InternalServerError() Response {
	return Response{
		Status: StatusError,
		Error:  "internal server error",
	}
}
