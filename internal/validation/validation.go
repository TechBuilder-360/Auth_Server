package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

func ValidateStruct(requestData interface{}, logger *log.Entry) ([]string, bool) {
	validationRes := validator.New()
	if err := validationRes.Struct(requestData); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		errMsgs := make([]string, 0)
		logger.Error("Validation failed on some fields : %+v", validationErrors)
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | validation failed '%s'",
				err.Field(),
				err.Value(),
				err.Tag(),
			))
		}

		return errMsgs, false
	}
	return nil, true
}
