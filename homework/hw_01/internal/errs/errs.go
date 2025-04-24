package errs

var (
	ErrIcomingDataIsEmpty = NewApiError(400, "incoming data is empty")
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

func (e *ApiError) Error() string {
	return e.Message
}
