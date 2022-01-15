package server

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type ErrorM map[string][]string

func (e ErrorM) Error() string {
	return "validation error"
}

func validationError(w http.ResponseWriter, _err error) {
	resp := ErrorM{}

	switch err := _err.(type) {
	case validator.ValidationErrors:
		for _, e := range err {
			field := e.Field()
			msg := checkTagRules(e)
			resp[field] = append(resp[field], msg)
		}
	default:
		resp["non_field_error"] = append(resp["non_field_error"], err.Error())
	}
	errorResponse(w, http.StatusUnprocessableEntity, resp)
}

func serverError(w http.ResponseWriter, err error) {
	log.Println(err)
	errorResponse(w, http.StatusInternalServerError, "internal error")
}

func errorResponse(w http.ResponseWriter, code int, errs interface{}) {
	writeJSON(w, code, M{"errors": errs})
}

func checkTagRules(e validator.FieldError) (errMsg string) {
	tag, field, param := e.ActualTag(), e.Field(), e.Param()

	if tag == "required" {
		errMsg = "this field is required"
	}

	if tag == "min" {
		errMsg = fmt.Sprintf("%s must be greater than %v", field, param)
	}

	if tag == "max" {
		errMsg = fmt.Sprintf("%s must be less than %v", field, param)
	}
	return
}
