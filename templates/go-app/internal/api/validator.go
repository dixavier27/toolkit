package api

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type nomePayload struct {
	Nome string `validate:"required,min=1,max=64"`
}

func validateNome(nome string) error {
	return validate.Struct(nomePayload{Nome: nome})
}
