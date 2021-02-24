package exchangeserve

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
const BASEURL = "https://api.exchangeratesapi.io/history?start_at=%s&end_at=%s&symbols=%s&base=%s"


/*
GetHistory returns a map of a decoded json object with
specified history of exchange rates based on date specified.
Optionally specify a certain country's currency or a base currency for comparison
*/
func GetHistory(startDate, endDate, currencyCode, currencyBase string) (map[string]interface{}, error) {
	var result = make(map[string]interface{})		// Body object
	// Insert parameters into BASEURL for request
	resData, err := http.Get(fmt.Sprintf(BASEURL, startDate, endDate, currencyCode, currencyBase))
	if err != nil { // Error handling data
		return nil, err
	}
	defer resData.Body.Close() // Closing body after finishing read
	if resData.StatusCode != 200 { // Error handling HTTP request
		e := errors.New(resData.Status)
		return nil, e
	}
	// Decoding body
	err = json.NewDecoder(resData.Body).Decode(&result)
	if err != nil { // Error handling decoding
		return nil, err
	}
	return result, err
}