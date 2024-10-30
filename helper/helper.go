package helper

import (
	"github.com/go-playground/validator/v10"
)

// Define the response structure
type response struct {
	Meta meta `json:"meta"`
	Data interface{} `json:"data"`
}

// Define the meta structure
type meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"` // This should be a single int representing an HTTP status code
	Status  string `json:"status"`
}

// APIResponse function returns a response object
func APIResponse(message string, code int, status string, data interface{}) response {

	// Construct the meta part of the response
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	// Construct and return the full response
	jsonResponse := response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}


func FormatValidationError(err error) []string {
	var errors []string;

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}


type UploadImageResponse struct {
	IsUploaded bool `json:"is_uploaded"`
}