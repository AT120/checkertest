package stdlib_helpers

import (
	"encoding/json"
	"io"
	"reflect"
)

type JsonConverter int

func (j JsonConverter) Decode(reader io.Reader, output any) error {
	return json.NewDecoder(reader).Decode(output)
}

func (j JsonConverter) LooselyCompare(a any, b any) (bool, error) {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, nil
	}

	switch a.(type) {
	case map[string]any:
		for key, bval := range b.(map[string]any) {
			aval, ok := a.(map[string]any)[key]
			if !ok {
				return false, nil
			}
			result, err := j.LooselyCompare(aval, bval)
			if !result || err != nil {
				return result, err
			}
		}
	case []any:
		used := make([]bool, len(a.([]any)))
		for _, bval := range b.([]any) {
			found := false
			for i, aval := range a.([]any) {
				if used[i] {
					continue
				}

				result, err := j.LooselyCompare(aval, bval)
				if err != nil {
					return result, err
				}

				if result {
					found = true
					used[i] = true
					break
				}
			}

			if !found {
				return false, nil
			}
		}

	case string:
		return TryRegexCompare(a.(string), b.(string))
	default:
		return a == b, nil
	}

	return true, nil
}

func (j JsonConverter) StrictlyCompare(a any, b any) (bool, error) {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, nil
	}

	switch a.(type) {
	case map[string]any:
		bmap := b.(map[string]any)
		amap := a.(map[string]any)
		if len(amap) != len(bmap) {
			return false, nil
		}
		for key, bval := range amap {
			aval, ok := amap[key]
			if !ok {
				return false, nil
			}
			result, err := j.StrictlyCompare(aval, bval)
			if !result || err != nil {
				return result, err
			}
		}

	case []any:
		aslice := a.([]any)
		bslice := b.([]any)
		if len(bslice) != len(aslice) {
			return false, nil
		}

		for i := range aslice {
			result, err := j.StrictlyCompare(aslice[i], bslice[i])
			if !result || err != nil {
				return result, err
			}
		}

	case string:
		return TryRegexCompare(a.(string), b.(string))
	default:
		return a == b, nil
	}

	return true, nil
}
