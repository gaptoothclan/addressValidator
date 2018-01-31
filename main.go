package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	add := Address{
		LineOne:   "Flat 12",    //", flat 5 69 sea road",
		LineTwo:   "Rose Tower", //"boscombe, bournemouth",
		LineThree: "62 Clarence Parade",
		Postcode:  "PO5 2HX",
	}

	ipc := IdealPostCode{
		ApiKey: "????",
	}

	//ap := FlatFileAddressProvider{}
	av := NewAddressValidator()
	addresses, err := av.ValidateAddress(add, ipc)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(len(addresses))
	for _, add := range addresses {
		b, _ := json.MarshalIndent(add, "", "    ")
		fmt.Println(string(b))
	}

}
