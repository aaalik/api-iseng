package validator

import (
	"fmt"
	"regexp"

	"github.com/aaalik/api-iseng/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

func init() {
	RegisterValidator("create_user", func(value interface{}) error {
		request, ok := value.(*model.RequestCreateUser)
		if !ok {
			return fmt.Errorf("invalid type, %T is not RequestCreateUser", value)
		}

		errs := validation.Errors{
			"name": validation.Validate(request.Name, validation.Required),
			"dob":  validation.Validate(request.Dob, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{4}[-][0-9]{2}[-][0-9]{2}$`))),
		}

		return errs.Filter()
	})

	RegisterValidator("update_user", func(value interface{}) error {
		request, ok := value.(*model.RequestUpdateUser)
		if !ok {
			return fmt.Errorf("invalid type, %T is not RequestCreateUser", value)
		}

		errs := validation.Errors{
			"id":   validation.Validate(request.ID, validation.Required),
			"name": validation.Validate(request.Name, validation.Required),
			"dob":  validation.Validate(request.Dob, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{4}[-][0-9]{2}[-][0-9]{2}$`))),
		}

		return errs.Filter()
	})
}
