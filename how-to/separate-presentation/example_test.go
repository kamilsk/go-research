package main

import (
	"encoding/json"
	"os"
)

func Example() {
	type Domain struct {
		Key string
	}

	type DomainInJSON struct {
		Key string `json:"key"`
	}

	converter := func(in struct {
		Key string
	}) DomainInJSON {
		return DomainInJSON{in.Key}
	}

	json.NewEncoder(os.Stdout).Encode(converter(Domain{Key: "value"}))

	// Output:
	// {"key":"value"}
}
