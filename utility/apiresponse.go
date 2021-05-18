package utility

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// ResponseObj ...
type ResponseObj struct {
	Status  bool
	Code    string
	Message string
}

// ResponseResultObj ...
type ResponseResultObj struct {
	ResponseObj
	Data interface{}
}

// ResponseResultObj ...
type ResponseValidateObj struct {
	ResponseObj
	ValidationErrors string
}

// NewResponse ...
func NewResponse() ResponseResultObj {
	return ResponseResultObj{}
}

// Success ...
func (res ResponseResultObj) PlainSuccess(code string, msg string) ResponseObj {

	response := ResponseObj{}
	response.Status = true
	response.Code = code
	response.Message = msg

	return response
}

// Success ...
func (res ResponseResultObj) Success(code string, msg string, data interface{}) ResponseResultObj {
	res.Status = true
	res.Code = code
	res.Message = msg
	if data.(interface{}) != "" {
		res.Data = data
	}

	fmt.Printf("response > %+v", res)
	return res
}

// Error ...
func (res ResponseResultObj) Error(code string, err string) ResponseObj {
	return ResponseObj{
		Status:  false,
		Code:    code,
		Message: err,
	}
}

func (res ResponseResultObj) ValidateError(code string, err string, errors validator.ValidationErrors) ResponseValidateObj {
	response := ResponseValidateObj{}
	response.Status = false
	response.Code = code
	response.Message = err
	response.ValidationErrors = errors.Error()
	return response
}
