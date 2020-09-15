package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestEqualFieldsSorted(t *testing.T) {
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

	assert.Equal(t, isEqual, false)
	assert.Equal(t, fieldError, "results.#.id")
}
