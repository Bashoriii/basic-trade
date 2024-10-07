package entity

import (
	"github.com/asaskevich/govalidator"
)

type Product struct {
	ID       int       `json:"id"`
	UUID     int       `json:"uuid"`
	Name     string    `json:"name" valid:"required~Product name is required"`
	ImageUrl string    `json:"image_url" valid:"url~Invalid image URL"`
	AdminID  int       `json:"admin_id"`
	Variants []Variant `json:"variants"`
}

func (a *Product) ValidateProduct() error {
	_, err := govalidator.ValidateStruct(a)
	if err != nil {
		return err
	}

	return nil
}
