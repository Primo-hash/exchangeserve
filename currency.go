package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/*
URL list for 'Foreign exchange rates API' to be modified to query needs
order: (startDate, endDate, currencyCode, currencyBase)
NOTE: 'symbols' and 'base' parameters can be omitted as query will return default values for these
 */
const HISTORYURL = "https://api.exchangeratesapi.io/history?start_at=%s&end_at=%s&symbols=%s&base=%s"
const CURRENCYURL = "https://api.exchangeratesapi.io/latest?&symbols=%s&base=%s"

/*
GetExchangeData returns a map of a decoded json object with
specified history of exchange rates based on date specified
or specify multiple countries' currencies and a base currency for comparison.
* Date parameters are mandatory when querying currency history.
* CurrencyCode is optional for history.
* CurrencyBase is optional.
*/
func GetExchangeData(startDate, endDate, currencyCode, currencyBase string) (map[string]interface{}, error) {
	if startDate == "" || endDate == "" { // Request for currency
		// Insert parameters into CURRENCYURL for request
		resData, err := http.Get(fmt.Sprintf(CURRENCYURL, currencyCode, currencyBase))
		if err != nil { // Error handling data
			return nil, err
		}
		return Decode(resData, "date")
	} else {							  // Request for history
		// Insert parameters into HISTORYURL for request
		resData, err := http.Get(fmt.Sprintf(HISTORYURL, startDate, endDate, currencyCode, currencyBase))
		if err != nil { // Error handling data
			return nil, err
		}
		return Decode(resData, "")
	}
}

// TODO make this function return a struct instead of a map, no random placement
// of JSON arrays and keys like map handles it

/*
Decode returns a decoded map from a decoded JSON
* Optional removal of a key in decoded map
*/
func Decode(data *http.Response, filter string) (map[string]interface{}, error) {
	var result = make(map[string]interface{})		// Body object

	defer data.Body.Close() // Closing body after finishing read
	if data.StatusCode != 200 { // Error handling HTTP request
		e := errors.New(data.Status)
		return nil, e
	}
	// Decoding body
	err := json.NewDecoder(data.Body).Decode(&result)
	if err != nil { // Error handling decoding
		return nil, err
	}
	// Optional filtering of certain key in map
	if filter != "" {
		// Check for filter word existence
		_, ok := result[filter];
		if ok {
			delete(result, filter);
		}
	}
	// Return map with requested data
	return result, err
}