package main

import (
	"encoding/json"
	"reflect"
	"strings"
)

const (
	diffType        = "diff-type"
	diffLengthBody  = "diff-length-body"
	diffLengthArray = "diff-length-array"
	diffValue       = "diff-value"
)

// Equal checks equality between 2 Body-encoded data.
func Equal(vx interface{}, vy interface{}) (bool, string) {

	if reflect.TypeOf(vx) != reflect.TypeOf(vy) {
		return false, diffType
	}

	switch x := vx.(type) {
	case map[string]interface{}:
		y := vy.(map[string]interface{})

		if len(x) != len(y) {
			return false, diffLengthBody
		}
		fieldError := ""
		for k, v := range x {
			val2 := y[k]

			if (v == nil) != (val2 == nil) {
				return false, k
			}

			isEqual, fieldErrorTemp := Equal(v, val2)
			if !isEqual {
				if fieldErrorTemp == diffValue || fieldErrorTemp == diffLengthArray || fieldErrorTemp == diffLengthBody {
					fieldErrorTemp = k
				} else {
					fieldErrorTemp = k + ".#." + fieldErrorTemp
				}
				if fieldError == "" {
					fieldError = fieldErrorTemp
				}
				return false, fieldError
			}
		}

		return true, "ok"
	case []interface{}:
		y := vy.([]interface{})

		if len(x) != len(y) {
			return false, diffLengthArray
		}

		fieldError := ""
		var matches int
		flagged := make([]bool, len(y))
		for _, v := range x {
			for i, v2 := range y {
				isEqual, fieldErrorTemp := Equal(v, v2)
				if isEqual && !flagged[i] {
					matches++
					flagged[i] = true
					break
				}else if !isEqual && fieldError == "" {
					fieldError = fieldErrorTemp
				}
			}
		}

		return matches == len(x), fieldError
	default:
		return vx == vy, diffValue
	}
}

func Remove(i interface{}, path string) {
	if path == "" {
		return
	}

	var next, current string
	split := strings.Split(path, ".#.")

	if len(split) == 1 {
		current = path
	} else if len(split) > 1 {
		current = split[0]
		next = strings.Join(split[1:], ".#.")
	}

	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			if k == current {
				// If there is no more nodes to traverse we can remove it and terminate the routine
				if next == "" {
					delete(t, current)
					return
				}
				Remove(v, next)
			}
		}
	case []interface{}:
		for _, v := range t {
			Remove(v, path)
		}
	}
}

// Unmarshal parses the Body-encoded data into an interface{}.
func Unmarshal(b []byte) (interface{}, error) {
	var j interface{}

	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	return j, nil
}
