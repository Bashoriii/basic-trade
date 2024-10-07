package entity

import "github.com/asaskevich/govalidator"

type Variant struct {
	ID          int    `json:"id"`
	UUID        int    `json:"uuid"`
	VariantName string `json:"variant_name" valid:"required~Variant name is required"`
	Quantity    int    `json:"quantity"`
	ProductID   int    `json:"product_id"`
}

func (a *Variant) ValidateVariant() error {
	_, err := govalidator.ValidateStruct(a)
	if err != nil {
		return err
	}

	return nil
}
