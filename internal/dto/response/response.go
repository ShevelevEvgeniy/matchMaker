package response

import "net/http"

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Status int `json:"status"`
}

func Created() SuccessResponse {
	return SuccessResponse{
		Status: http.StatusCreated,
	}
}

func BadRequest(msg string) Response {
	return Response{
		Status: http.StatusBadRequest,
		Error:  msg,
	}
}

func InternalServerError() Response {
	return Response{
		Status: http.StatusInternalServerError,
		Error:  "internal server error",
	}
}
