package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type AuthForm struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required"`
}

func (u *AuthForm) Hash256Password(password string) string {
	buf := []byte(password)
	pwd := sha3.New256()
	pwd.Write(buf)
	return hex.EncodeToString(pwd.Sum(nil))
}

func (u *AuthForm) BcryptPassword(password string) (string, error) {
	buf := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(buf, bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

func (u *AuthForm) ValidateBcryptPassword(password, password2 string) bool {
	byteHash := []byte(password)
	buf := []byte(password2)
	err := bcrypt.CompareHashAndPassword(byteHash, buf)
	if err != nil {
		return false
	}
	return true
}

func (u *AuthForm) Validate(v *validator.Validate) error {
	err := v.Struct(u)
	if err != nil {

		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			var out error
			err = err.(validator.ValidationErrors)
			for _, fe := range ve {
				out = errors.New(fmt.Sprintf("%s. %s", fe.StructField(), msgForTag(fe.Tag())))
			}
			return out
		}
	}
	return nil
}
