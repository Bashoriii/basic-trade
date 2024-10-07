package entity

import (
	"basic-trade/helpers"

	"github.com/asaskevich/govalidator"
)

type Admin struct {
	ID       int       `json:"id"`
	UUID     int       `json:"uuid"`
	Name     string    `json:"name" valid:"required~Your name is required"`
	Email    string    `json:"email" valid:"required~Your email is required, email~Invalid email format"`
	Password string    `json:"password" valid:"required~Your password is required, minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Products []Product `json:"products"`
}

func (a *Admin) ValidateAdmin() error {
	_, err := govalidator.ValidateStruct(a)
	if err != nil {
		return err
	}

	a.Password = helpers.HashingPassword(a.Password)
	return nil
}
