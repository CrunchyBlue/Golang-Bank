package util

import "github.com/CrunchyBlue/Golang-Bank/constants"

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case constants.USD, constants.EUR, constants.CAD:
		return true
	}
	return false
}
