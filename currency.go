package exchangeserve



/*
GetHistory returns a map of a decoded json object with
specified country's history of exchange rates based on date specified.
*/
func GetCurrenc(countryName, datePeriod string) ()
resp, err := http.Get("http://example.com/")