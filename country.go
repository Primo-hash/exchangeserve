package exchange

import (
	"github.com/alediaferia/gocountries"
)

/*
GetCurrency returns a string of specified Country's currency code e.g.(NOK, USD, EUR...)
*/
func GetCurrency(countryName string) (string, error) {
	// Query for structs of possible countries
	countries, err := gocountries.CountriesByName(countryName)
	// Extract first country
	country := (countries)[0]
	// Extract currency code
	currencyCode := country.Currencies[0]
	return currencyCode, err
}
