package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestEqual(t *testing.T) {
	jsonFileArrayBytes1, _ := ioutil.ReadFile("response-host-1.json")
	jsonFileArrayBytes2, _ := ioutil.ReadFile("response-host-2.json")

	leftJSON, _ := unmarshal(jsonFileArrayBytes1)
	rightJSON, _ := unmarshal(jsonFileArrayBytes2)

	excludeFields := []string{"paging", "results.#.payer_costs.#.payment_method_option_id"}

	if len(excludeFields) > 0 {
		for _, excludeField := range excludeFields {
			Remove(leftJSON, excludeField)
			Remove(rightJSON, excludeField)
		}
	}

	isEqual, fieldError := Equal(leftJSON, rightJSON)

	fmt.Println("isEqual:", isEqual, ", FieldError:", fieldError)
}
