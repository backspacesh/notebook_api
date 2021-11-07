package model

import validation "github.com/go-ozzo/ozzo-validation"

func requiredIf(cond bool) validation.RuleFunc {
	return func(data interface{}) error {
		if cond {
			validation.Validate(data, validation.Required)
		}

		return nil
	}
}