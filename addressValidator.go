package main

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

//
type AddressProvider interface {
	GetAddressData(postcode string) ([]Address, error)
}

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
	Address       *Address
	Score         int
	ClosePenalty  int
	Tokens        []string
	PrimaryTokens []string
	Matches       []string
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
func (av *AddressValidator) ValidateAddress(address Address, addressProvider AddressProvider) ([]TokenisedAddress, error) {
	addresses, err := addressProvider.GetAddressData(address.Postcode)
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		tokenisedAddress := av.tokeniseAddress(address)
		av.tokenisedAddress = append(av.tokenisedAddress, tokenisedAddress)
	}

	rankedAddresses := av.rankAddresses(address)

	// only one address or no addresses
	if len(rankedAddresses) < 2 {
		return rankedAddresses, nil
	}

	// remove all but the highest ranked
	highestScore := rankedAddresses[0].Score
	var filteredAddresses []TokenisedAddress
	for _, add := range rankedAddresses {
		if add.Score == highestScore {
			filteredAddresses = append(filteredAddresses, add)
		}
	}

	if len(filteredAddresses) == 1 {
		return filteredAddresses, nil
	}

	lowestClosePenalty := rankedAddresses[0].ClosePenalty
	if lowestClosePenalty < rankedAddresses[1].ClosePenalty {
		return []TokenisedAddress{rankedAddresses[0]}, nil
	}

	return []TokenisedAddress{}, nil
}

func (av *AddressValidator) rankAddresses(address Address) []TokenisedAddress {
	tokenisedCheck := av.tokeniseAddress(address)

	for i := range av.tokenisedAddress {
		tokenCount := len(av.tokenisedAddress[i].Tokens)
		checkAddress := &av.tokenisedAddress[i]
		matchCount := 0
		score := 0
		for _, token := range checkAddress.Tokens {

			if av.inArray(tokenisedCheck.Tokens, token) {
				matchCount++
				checkAddress.Matches = append(checkAddress.Matches, token)
				if av.inArray(checkAddress.PrimaryTokens, token) {
					score = score + 2
				} else {
					score++
				}
			}
		}

		checkAddress.ClosePenalty = int(math.Abs(float64(matchCount - tokenCount)))
		checkAddress.Score = score // maybe a percentage
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

func (av *AddressValidator) inArray(haystack []string, needle string) bool {
	for _, item := range haystack {
		if needle == item {
			return true
		}
	}
	return false
}

func (av *AddressValidator) tokeniseAddress(address Address) TokenisedAddress {
	combined := strings.Join([]string{
		address.LineOne,
		address.LineTwo,
		address.LineThree}, " ")

	buildingName := av.splitString(address.BuildingName)
	buildingNumber := av.splitString(address.BuildingNumber)
	subBuildingName := av.splitString(address.SubBuildingName)

	primaryTokens := append(buildingName, buildingNumber...)
	primaryTokens = append(primaryTokens, subBuildingName...)

	filteredTokens := av.splitString(combined)

	return TokenisedAddress{
		Address:       &address,
		Tokens:        filteredTokens,
		PrimaryTokens: primaryTokens,
	}
}

func (av *AddressValidator) splitString(toSplit string) []string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	combined := reg.ReplaceAllString(toSplit, "")

	combined = strings.ToLower(combined)

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

	return filteredTokens
}
