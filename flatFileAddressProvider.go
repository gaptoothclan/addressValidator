package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

//
type FlatFileAddressProvider struct{}

//
func (f FlatFileAddressProvider) GetAddressData(postcode string) ([]Address, error) {
	postcode = strings.ToLower(postcode)
	postcode = strings.Replace(postcode, " ", "", -1)
	filename := fmt.Sprintf("./data/%s.json", postcode)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	addressResult := AddressResult{}
	err = json.Unmarshal(file, &addressResult)
	if err != nil {
		return nil, err
	}

	return addressResult.Result, nil
}
