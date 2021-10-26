package main

import (
	"context"
	"github.com/aliparlakci/armut-backend-assessment/services"
	"github.com/go-playground/validator/v10"
)

func ExistingUserValidator(getter services.UserGetter) validator.Func {
	return func(fl validator.FieldLevel) bool {
		username, ok := fl.Field().Interface().(string)
		if ok {
			exists, err := getter.UserExists(context.TODO(), username)
			if err != nil {
				// TODO: Log
			}

			return exists
		}

		return false
	}
}
