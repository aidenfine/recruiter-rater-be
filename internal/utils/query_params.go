package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type MissingParamError struct {
	Message string
}

func (e *MissingParamError) Error() string {
	return e.Message
}

func ParseQueryParams[T any](r *http.Request) (*T, error) {
	query := r.URL.Query()
	var res T
	val := reflect.ValueOf(&res).Elem() // allows us to read key of struct ex T.Name, T.Limit
	typ := val.Type()                   // allows us to read fields

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		queryKey := field.Tag.Get("query")

		// default to name if no query key found
		// TODO: may want to raise an error of some kind to prevent this from happening
		if queryKey == "" {
			queryKey = field.Name
		}
		// check if field is required
		isRequired := field.Tag.Get("required") == "true"

		// get val and return err if empty and required
		value := query.Get(queryKey)
		if isRequired && value == "" {
			return nil, &MissingParamError{Message: fmt.Sprintf("Missing query param: %s", queryKey)}
		}

		if value == "" {
			continue
		}

		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		// create error message based on type defined in struct
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(value)
		case reflect.Int:
			intVal, err := strconv.Atoi(value)
			if err != nil {
				return nil, &MissingParamError{Message: fmt.Sprintf("Invalid int for %s: %s", queryKey, value)}
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return nil, &MissingParamError{Message: fmt.Sprintf("Invalid bool for %s: %s", queryKey, value)}
			}
			fieldVal.SetBool(boolVal)
		default:
			return nil, &MissingParamError{Message: fmt.Sprintf("Unsupported field type for %s", field.Name)}
		}
	}
	return &res, nil
}
