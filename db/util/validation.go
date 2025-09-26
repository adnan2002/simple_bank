package util

import "github.com/go-playground/validator/v10"




var Currency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return validCurrency(currency)
	}
	return false
}







func validCurrency(curr string) bool {
	switch curr {
		case "EUR", "BHD", "USD", "AED", "SAR", "CAD" :
			return true
		default:
			return false
	}	
}