package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var statusCode = map[string]int{
	"Bad Request":           http.StatusBadRequest,
	"Unauthorized":          http.StatusUnauthorized,
	"Forbidden":             http.StatusForbidden,
	"Not Found":             http.StatusNotFound,
	"Method Not Allowed":    http.StatusMethodNotAllowed,
	"Internal Server Error": http.StatusInternalServerError,
}

// Error structure
type Error struct {
	ID      string
	params  map[string]interface{}
	message string // Translated error message
}

// NewError creates a new error
func NewError(id string, params map[string]interface{}) *Error {
	er := Error{
		ID:     id,
		params: params,
	}
	return &er
}

// MarshalJSON : custom json marshal func for MyUUID
func (er *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(er.message)
}

// AppError structure
type AppError struct {
	// for the end user
	Errors     []map[string]interface{}
	ErrorCode  string `json:"error"`
	Message    string `json:"detail"` // Message to be display to the end user without debugging information
	StatusCode int    `json:"-"`
	Many       bool   `json:"-"`

	// for internal debug
	Where         string `json:"-"` // The function where it happened in the form of Struct.Func
	DetailedError string `json:"-"` // Internal error string to help the developer
	ClientID      string `json:"-"` // The ClientID that's also set in the header
}

func (er *AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

// NewAppError creates a new app error
func NewAppError(where string, id string, params map[string]interface{}, details string, status int) *AppError {
	ap := &AppError{}
	ap.Errors = []map[string]interface{}{{"detail": []*Error{{ID: id, params: params}}}}
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = status
	ap.ErrorCode = "error"
	return ap
}

// NewAppErrorWithDetails creates a new app error
func NewAppErrorWithDetails(where string, errors []map[string]interface{}, many bool, code string, status int) *AppError {
	ap := &AppError{}
	ap.Message = ""
	ap.Errors = errors
	ap.Many = many
	ap.Where = where
	ap.DetailedError = ""
	ap.StatusCode = status
	ap.ErrorCode = code
	return ap
}

// NewAppErrorWithCode creates a new app error with status code
func NewAppErrorWithCode(where string, id string, params map[string]interface{}, details, code string, status int) *AppError {
	ap := &AppError{}
	ap.Errors = []map[string]interface{}{{"detail": []*Error{{ID: id, params: params}}}}
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = status
	ap.ErrorCode = code
	return ap
}

// NewSimpleAppError creates a new simple app error with status code
func NewSimpleAppError(where string, id string, params map[string]interface{}, status int) *AppError {
	return NewAppErrorWithCode(where, id, params, "", "", status)
}

// SystemMessage translates app error to system message
func (er *AppError) SystemMessage() string {
	if len(er.Errors) == 1 {
		if err, ok := er.Errors[0]["detail"]; ok {
			errAry := err.([]*Error)
			return errAry[0].message
		}
	}
	data := []map[string]interface{}{}
	for _, errMap := range er.Errors {
		temp := map[string]interface{}{}
		for k, v := range errMap {
			b, _ := json.Marshal(v)
			temp[k] = string(b)
		}
		data = append(data, temp)
	}

	if len(data) == 1 {
		b, _ := json.Marshal(data[0])
		return string(b)
	}

	b, _ := json.Marshal(data)
	return string(b)
}

// ToJSON convert a AppError to a json string
func (er *AppError) ToJSON() string {
	// customize error response
	r := make(map[string]interface{})

	data := []map[string]interface{}{}
	for _, errMap := range er.Errors {
		temp := map[string]interface{}{}
		for k, v := range errMap {
			if k == "detail" {
				errAry := v.([]*Error)
				temp[k] = errAry[0]
			} else {
				temp[k] = v
			}
		}
		data = append(data, temp)
	}
	if er.ErrorCode == "" {
		r["error"] = data[0]["detail"]
	} else {
		if er.ErrorCode == "error" {
			r["error"] = http.StatusText(er.StatusCode)
		} else {
			r["error"] = er.ErrorCode
		}
		if !er.Many {
			r["data"] = data[0]
		} else {
			r["data"] = data
		}
	}

	b, _ := json.Marshal(r)
	return string(b)
}

// AppErrorFromJSON will decode the input and return an AppError
func AppErrorFromJSON(data io.Reader) *AppError {
	var er struct {
		Data struct {
			Detail string `json:"detail"`
		} `json:"data"`
		Error string `json:"error"`
	}

	err := json.NewDecoder(data).Decode(&er)
	if err == nil {
		return NewAppError("", "", nil, er.Data.Detail, statusCode[er.Error])
	}
	return NewAppError("AppErrorFromJSON", "model.utils.decode_json.app_error", nil, err.Error(), http.StatusInternalServerError)
}

// InvalidFieldError creates new invalid field error
func InvalidFieldError(tableName, fieldName string, id interface{}) *AppError {
	msgID := fmt.Sprintf("model.%s.is_valid.%s.app_error", tableName, fieldName)
	details := ""
	if id != nil {
		details = fmt.Sprintf("id=%v", id)
	}
	return ValidationError("InvalidFieldError", msgID, nil, details)
}

// ValidationError creates new validation error
func ValidationError(where, message string, params map[string]interface{}, details string) *AppError {
	if len(message) == 0 {
		message = "Invalid input"
	}
	return NewAppErrorWithCode(where, message, params, details, "ValidationError", http.StatusBadRequest)
}

// ValidationErrorWithDetails creates new validation error
func ValidationErrorWithDetails(where string, errors []map[string]interface{}) *AppError {
	return NewAppErrorWithDetails(where, errors, false, "ValidationError", http.StatusBadRequest)
}

// ValidationErrorWithManyDetails creates new validation error
func ValidationErrorWithManyDetails(where string, errors []map[string]interface{}) *AppError {
	return NewAppErrorWithDetails(where, errors, true, "ValidationError", http.StatusBadRequest)
}

// ParseError creates new parse error
func ParseError(where, id string, params map[string]interface{}, details string) *AppError {
	return NewAppErrorWithCode(where, id, nil, details, "ParseError", http.StatusBadRequest)
}

// UnauthorizedError creates new authorization failed error
func UnauthorizedError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.unauthorized", nil, details, "Unauthorized", http.StatusUnauthorized)
}

// AuthenticationFailedError creates new authentication failed error
func AuthenticationFailedError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.authentication_failed", nil, details, "AuthenticationFailed", http.StatusUnauthorized)
}

// AuthenticationFailedCustomError creates new custom authentication failed error
func AuthenticationFailedCustomError(where, message string, params map[string]interface{}, details string) *AppError {
	return NewAppErrorWithCode(where, message, params, details, "AuthenticationFailed", http.StatusUnauthorized)
}

// NotAuthenticatedError creates new not authenticated error
func NotAuthenticatedError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.not_authenticated", nil, details, "NotAuthenticated", http.StatusUnauthorized)
}

// PermissionDeniedError creates new permission denied error
func PermissionDeniedError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.permission_denied", nil, details, "PermissionDenied", http.StatusForbidden)
}

// NotFoundError creates new not found error
func NotFoundError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.not_found", nil, details, "NotFound", http.StatusNotFound)
}

// ExpiredTokenError creates new token expired error
func ExpiredTokenError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.expired_token", nil, details, "ExpiredTokenError", http.StatusBadRequest)
}

// InvalidTokenError creates new invalid token error
func InvalidTokenError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.invalid_token", nil, details, "InvalidTokenError", http.StatusBadRequest)
}

// InvalidAuthenticationTokenError creates new invalid token error
func InvalidAuthenticationTokenError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.invalid_authentication_token", nil, details, "InvalidAuthenticationTokenError", http.StatusBadRequest)
}

// ExpiredAuthenticationTokenError creates new invalid token error
func ExpiredAuthenticationTokenError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.expired_authentication_token", nil, details, "InvalidAuthenticationTokenError", http.StatusBadRequest)
}

// InvalidAuthenticationTokenError creates new invalid token error
func InvalidSmsUpdateTokenError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.invalid_token", nil, details, "InvalidSmsUpdateTokenError", http.StatusBadRequest)
}

// InvalidAuthenticationTokenError creates new invalid token error
func InvalidEmailUpdateTokenError(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.invalid_token", nil, details, "InvalidEmailUpdateTokenError", http.StatusBadRequest)
}

// LoginUserNotFound creates new user not found error
func LoginUserNotFound(where, details string) *AppError {
	return NewAppErrorWithCode(where, "model.app_error.login_user_not_found", nil, details, "LoginUserNotFound", http.StatusBadRequest)
}

// ParamNotFoundError creates new param not found error
func ParamNotFoundError(where, parameter string) *AppError {
	return NewSimpleAppError(where, "api.context.param.not_found", map[string]interface{}{"Param": parameter}, http.StatusBadRequest)
}

// InvalidParamError creates new invalid request param error
func InvalidParamError(parameter string) *AppError {
	return ValidationError("", "api.context.invalid_body_param", map[string]interface{}{"Name": parameter}, "")
}

// InvalidRequestBody creates new invalid request param error
func InvalidRequestBodyError() *AppError {
	return ValidationError("", "api.context.invalid_request_body", nil, "")
}

// InvalidUrlParamError creates new invalid url param error
func NewInvalidUrlParamError(parameter string) *AppError {
	return ValidationError("", "api.context.invalid_url_param", map[string]interface{}{"Name": parameter}, "")
}

// SameCurrentNoteError creates new same current note error
func SameCurrentNoteError(where, details string) *AppError {
	return ValidationError(where, "model.app_error.same_current_note", nil, details)
}

// SameCurrentTaskError create new same current task error
func SameCurrentTaskError(where, details string) *AppError {
	return ValidationError(where, "model.app_error.same_current_task", nil, details)
}
