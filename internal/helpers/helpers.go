package helpers

import (
	"math"
	"strings"
	"unicode"
)

func NormalizeZipCode(zipcode string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}

		return -1 // -1 descarta o caractere
	}, zipcode)
}

func ValidateZipCode(zipcode string) bool {
	normalized := NormalizeZipCode(zipcode)
	
	return len(normalized) == 8
}

func CelsiusToFahrenheit(celsius float64) float64 {
	return math.Ceil((celsius * 9 / 5) + 32)
}

func CelsiusToKelvin(celsius float64) float64 {
	return math.Ceil(celsius + 273.15)
}
