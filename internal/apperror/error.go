package apperror

import "encoding/json"

var ErrNotFound = NewAppErr(nil, "not found", "", "US-000003")

type AppError struct {
	Err              error  `json:"_"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

// for implementing error interface
func (e *AppError) Error() string {
	return e.Message

}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}
func NewAppErr(err error, message string, developerMessage string, code string) *AppError {
	return &AppError{
		Err:              err, // fmt.Errorf(message),
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(err error) *AppError {
	return NewAppErr(err, "internal system error", err.Error(), "US-000000")
}
