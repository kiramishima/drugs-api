package models

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type DrugForm struct {
	Name        *string `json:"name" db:"name" validate:"required"`
	Approved    *bool   `json:"approved" db:"approved" validate:"required"`
	MinDose     *int    `json:"min_dose" db:"min_dose" validate:"required"`
	MaxDose     *int    `json:"max_dose" db:"max_dose" validate:"required"`
	AvailableAt *string `json:"available_at" db:"available_at" validate:"required"`
}

func (u *DrugForm) Validate(v *validator.Validate) error {
	err := v.Struct(u)
	if err != nil {

		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			var out error
			err = err.(validator.ValidationErrors)
			for _, fe := range ve {
				out = errors.New(fmt.Sprintf("%s: %s", fe.Field(), msgForTag(fe.Tag())))
			}
			return out
		}
	}
	return nil
}
