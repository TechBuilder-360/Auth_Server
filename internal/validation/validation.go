package validation

import (
	"errors"
	"fmt"
	"github.com/TechBuilder-360/Auth_Server/pkg/log"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func ValidateStruct(requestData interface{}, logger log.Entry) (string, bool) {
	validationRes := validator.New()
	if err := validationRes.Struct(requestData); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		errArr := make([]string, 0)
		logger.Error("Validation failed on some fields : %+v", validationErrors)
		for _, e := range validationErrors {
			fieldName := e.Field()
			field, _ := reflect.TypeOf(requestData).Elem().FieldByName(fieldName)
			fieldJSONName, _ := field.Tag.Lookup("json")

			errArr = append(errArr, fmt.Sprintf(
				"[%s]--> '%v' | validation failed '%s'",
				fieldJSONName,
				e.Value(),
				e.Tag(),
			))

		}

		return strings.Join(errArr, "\n"), false
	}
	return "", true
}
