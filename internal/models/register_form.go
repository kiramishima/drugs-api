package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

// RegisterForm struct para la data del registro
type RegisterForm struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password,omitempty" validate:"required"`
	Name     string `form:"name" json:"name,omitempty"`
}

func (u *RegisterForm) Hash256Password(password string) string {
	buf := []byte(password)
	pwd := sha3.New256()
	pwd.Write(buf)
	return hex.EncodeToString(pwd.Sum(nil))
}

func (u *RegisterForm) BcryptPassword(password string) (string, error) {
	buf := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(buf, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

func (u *RegisterForm) ValidateBcryptPassword(password, password2 string) bool {
	byteHash := []byte(password)
	buf := []byte(password2)
	err := bcrypt.CompareHashAndPassword(byteHash, buf)
	if err != nil {
		return false
	}
	return true
}

func (u *RegisterForm) Validate(v *validator.Validate) error {
	err := v.Struct(u)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var out error
		err = err.(validator.ValidationErrors)
		for _, fe := range ve {
			out = errors.New(fmt.Sprintf("%s. %s", fe.StructField(), msgForTag(fe.Tag())))
		}
		return out
	}
	return nil
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Bad email format"
	case "gt":
		return "The value needs to be more than 0 and non negative"
	}
	return ""
}
