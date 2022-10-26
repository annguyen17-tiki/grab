package handler

import (
	"regexp"

	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/go-playground/validator/v10"
)

var defaultPagination = model.Paging{
	Limit:    10,
	Offset:   0,
	OrderBy:  "created_at",
	OrderDir: "DESC",
}

var globalValidator *validator.Validate

func initValidator() error {
	v := validator.New()
	if err := v.RegisterValidation("phone", validatePhone); err != nil {
		return err
	}

	globalValidator = v
	return nil
}

func validatePhone(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^(0|84|084|\+84)[35789](\d{8})$`).MatchString(fl.Field().String())
}
