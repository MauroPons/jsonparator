package main

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

const (
	diffType        = "diff-type"
	diffLengthBody  = "diff-length-body"
	diffLengthArray = "diff-length-array"
	diffValue       = "diff-value"
	diffArrayEmpty  = "diff-array-empty"
	diffBodyEmpty   = "diff-body-empty"
)

//func Equal(vx interface{}, vy interface{}) (bool, string) {
//
//	if reflect.TypeOf(vx) != reflect.TypeOf(vy) {
//		return false, diffType
//	}
//	switch x := vx.(type) {
//	case map[string]interface{}:
//		y := vy.(map[string]interface{})
//		lengthMapX := len(x)
//		lengthMapY := len(y)
//
//		if lengthMapX > 0 && lengthMapY == 0 {
//			return false, diffBodyEmpty
//		}
//
//		if lengthMapX != lengthMapY {
//			if lengthMapX > lengthMapY {
//				for k, _ := range y {
//					delete(x, k)
//				}
//				var arrayKeys []string
//				for k := range x {
//					arrayKeys = append(arrayKeys, k)
//				}
//				sort.Strings(arrayKeys)
//				fieldsDiff := ".#.sc1.#." + strings.Join(arrayKeys[:], ".#.")
//				return false, diffLengthBody + fieldsDiff
//			} else if lengthMapY > lengthMapX {
//				for k, _ := range x {
//					delete(y, k)
//				}
//				var arrayKeys []string
//				for k := range y {
//					arrayKeys = append(arrayKeys, k)
//				}
//				sort.Strings(arrayKeys)
//				fieldsDiff := ".#.sc2.#." + strings.Join(arrayKeys[:], ".#.")
//				return false, diffLengthBody + fieldsDiff
//			}
//			return false, diffLengthBody
//		}
//		var arrayKeys []string
//		for k := range x {
//			arrayKeys = append(arrayKeys, k)
//		}
//		sort.Strings(arrayKeys)
//
//		for i := range arrayKeys {
//			k := arrayKeys[i]
//			xv := x[k]
//			yv := y[k]
//			isEqual, fieldError := Equal(xv, yv)
//			if !isEqual {
//				return false, k + ".#." + fieldError
//			}
//		}
//		return true, ""
//	case []interface{}:
//		y := vy.([]interface{})
//		if len(x) == 0 && len(y) == 0 {
//			return true, ""
//		}
//		if len(x) > 0 && len(y) == 0 {
//			return false, diffArrayEmpty
//		}
//		if len(x) != len(y) {
//			return false, diffLengthArray
//		}
//
//		switch x[0].(type) {
//		case string:
//			var xString []string
//			var yString []string
//			for _, param := range x {
//				xString = append(xString, param.(string))
//			}
//			for _, param := range y {
//				yString = append(yString, param.(string))
//			}
//			sort.Strings(xString)
//			sort.Strings(yString)
//
//			for index := range xString {
//				isEqual, fieldError := Equal(xString[index], yString[index])
//				if !isEqual {
//					return false, fieldError
//				}
//			}
//		default:
//			for index := range x {
//				isEqual, fieldError := Equal(x[index], y[index])
//				if !isEqual {
//					return false, fieldError
//				}
//			}
//		}
//
//		return true, ""
//	default:
//		areEquals := vx == vy
//		var fieldError string
//		if !areEquals {
//			fieldError = diffValue
//		}
//		return areEquals, fieldError
//	}
//}

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
		var arrayKeys []string
		for k := range x {
			arrayKeys = append(arrayKeys, k)
		}
		sort.Strings(arrayKeys)

		for i := range arrayKeys {
			k := arrayKeys[i]
			xv := x[k]
			yv := y[k]
			isEqual, fieldError := Equal(xv, yv)
			if !isEqual {
				if fieldError == diffLengthBody || fieldError == diffLengthArray || fieldError == diffType || fieldError == diffValue {
					return false, k
				} else {
					return false, k + ".#." + fieldError
				}
				//return false, k + ".#." + fieldError
			}
		}
		return true, ""
	case []interface{}:
		y := vy.([]interface{})
		if len(x) != len(y) {
			return false, diffLengthArray
		}
		for index := range x {
			isEqual, fieldError := Equal(x[index], y[index])
			if !isEqual {
				return false, fieldError
			}
		}
		return true, ""
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
