package models

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type VaccinationForm struct {
	Name      *string `json:"name" db:"name" validate:"required"`
	DrugID    *int    `json:"drug_id" db:"drug_id" validate:"required"`
	Dose      *int    `json:"dose" db:"dose" validate:"required"`
	AppliedAt *string `json:"applied_at" db:"applied_at" validate:"required"`
}

func (u *VaccinationForm) Validate(v *validator.Validate) error {
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
