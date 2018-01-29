package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strings"
)

// AddressResult Holds an address result from ideal postcodes
type AddressResult struct {
	Result  []Address `json:"result"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
}

// Address data type, partial data from ideal postcodes
type Address struct {
	LineOne         string `json:"line_1"`
	LineTwo         string `json:"line_2"`
	LineThree       string `json:"line_3"`
	BuildingNumber  string `json:"building_number"`
	BuildingName    string `json:"building_name"`
	SubBuildingName string `json:"sub_building_name"`
	Postcode        string `json:"postcode"`
}

// TokenisedAddress holds the address and the scoring
type TokenisedAddress struct {
	Address      *Address
	Score        int
	ClosePenalty int
	Tokens       []string
	Matches      []string
}

// AddressValidator the main struct for validating addresses
type AddressValidator struct {
	tokenisedAddress []TokenisedAddress
}

// NewAddressValidator return an AddressValidator
func NewAddressValidator() *AddressValidator {
	return &AddressValidator{}
}

// ValidateAddress entry point
func (av *AddressValidator) ValidateAddress(address Address) ([]TokenisedAddress, error) {
	addresses, err := av.getAddressData(address.Postcode)
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		tokenisedAddress := av.tokeniseAddress(address)
		av.tokenisedAddress = append(av.tokenisedAddress, tokenisedAddress)
	}

	rankedAddresses := av.rankAddresses(address)

	return rankedAddresses, nil
}

func (av *AddressValidator) rankAddresses(address Address) []TokenisedAddress {
	tokenisedCheck := av.tokeniseAddress(address)

	for i := range av.tokenisedAddress {
		tokenCount := len(av.tokenisedAddress[i].Tokens)
		checkAddress := &av.tokenisedAddress[i]
		// What if one has more items than the other
		// compare each item in address and check address
		// most matches wins
		// Maybe record the matched items
		// percentage score
		// need a filtering system to break up things
		// like 47a
		matchCount := 0
		for _, token := range checkAddress.Tokens {
			for j := range tokenisedCheck.Tokens {
				if tokenisedCheck.Tokens[j] == token {
					matchCount++
					checkAddress.Matches = append(checkAddress.Matches, token)
				}
			}
		}

		checkAddress.ClosePenalty = int(math.Abs(float64(matchCount - tokenCount)))
		checkAddress.Score = matchCount // maybe a percentage
	}

	// Sort by score then by close penalty
	sort.Slice(av.tokenisedAddress, func(i, j int) bool {
		if av.tokenisedAddress[i].Score > av.tokenisedAddress[j].Score {
			return true
		}
		if av.tokenisedAddress[i].Score < av.tokenisedAddress[j].Score {
			return false
		}
		return av.tokenisedAddress[i].ClosePenalty < av.tokenisedAddress[j].ClosePenalty
	})

	return av.tokenisedAddress
}

func (av *AddressValidator) getAddressData(postcode string) ([]Address, error) {
	postcode = strings.ToLower(postcode)
	postcode = strings.Replace(postcode, " ", "", -1)
	filename := fmt.Sprintf("./%s.json", postcode)

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

func (av *AddressValidator) tokeniseAddress(address Address) TokenisedAddress {
	combined := strings.Join([]string{
		address.LineOne,
		address.LineTwo,
		address.LineThree}, " ")

	combined = strings.ToLower(combined)

	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	combined = reg.ReplaceAllString(combined, "")

	tokens := strings.Split(combined, " ")

	snReg, _ := regexp.Compile(`^([0-9]+)([a-z]+)$`)

	var filteredTokens []string
	for _, token := range tokens {
		snMatch := snReg.FindStringSubmatch(token)

		if len(snMatch) == 3 {
			for _, match := range snMatch[1:] {
				filteredTokens = append(filteredTokens, match)
			}
		} else if token != "" {
			filteredTokens = append(filteredTokens, token)
		}
	}

	sort.Strings(filteredTokens)

	return TokenisedAddress{
		Address: &address,
		Tokens:  filteredTokens,
	}
}

/*
add weightings these should be higher
5's
building_number
building_name
sub_building_name 12a both components are important

thoroughfare

*/
func main() {
	add := Address{
		LineOne:   "Flat 20",    //", flat 5 69 sea road",
		LineTwo:   "Rose Tower", //"boscombe, bournemouth",
		LineThree: "62 clarence parade",
		Postcode:  "PO5 2HX",
	}

	av := NewAddressValidator()
	addresses, err := av.ValidateAddress(add)

	if err != nil {
		fmt.Print(err)
		return
	}

	for _, add := range addresses {
		fmt.Printf("%s\t%s\t%s\t%d\t%d\t%d\n", add.Address.LineOne, add.Address.LineTwo, add.Address.LineThree, add.Score, len(add.Tokens), add.ClosePenalty)
	}

}
