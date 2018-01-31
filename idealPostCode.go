package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const IDEAL_POSTCODE_URL = "https://api.ideal-postcodes.co.uk/v1/postcodes/%s?api_key=%s"

//
type IdealPostCode struct {
	ApiKey string
}

//
func (i IdealPostCode) GetAddressData(postcode string) ([]Address, error) {
	url := fmt.Sprintf(IDEAL_POSTCODE_URL, postcode, i.ApiKey)

	response, err := http.Get(url)

	defer response.Body.Close()

	if err != nil {
		return nil, err
	}

	contents, _ := ioutil.ReadAll(response.Body)

	var result AddressResult
	jsonErr := json.Unmarshal(contents, &result)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return result.Result, nil
}
