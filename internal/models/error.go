package models

type ApiError struct {
	Error string       `json:"error"`
	Data  ApiErrorData `json:"data"`
}

type ApiErrorData struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

func ErrorToApiError(err error) *ApiError {
	return &ApiError{
		Error: err.Error(),
		Data: ApiErrorData{
			Error:   err.Error(),
			Details: err.Error(),
		},
	}
}
