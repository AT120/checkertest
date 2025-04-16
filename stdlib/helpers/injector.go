package stdlib_helpers

import (
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"fmt"
	"strconv"
	"strings"
)

type stringIterator struct {
	str          string
	start        string
	stop         string
	i            int
	prefix_end   int
	prefix_begin int
}

func (si *stringIterator) End() bool {
	return si.prefix_end >= len(si.str)
}

func (si *stringIterator) Next() (string, bool) {
	si.prefix_begin = si.i
	if si.i > 0 {
		si.prefix_begin += len(si.stop)
	}

	var (
		startIdx int = -1
		endIdx   int = -1
	)

	i := si.i

	for start_i := 0; i < len(si.str); i++ {
		if si.str[i] == si.start[start_i] {
			start_i++
			if start_i == len(si.start) {
				si.prefix_end = i - start_i + 1
				startIdx = i + 1
				break
			}
		} else {
			start_i = 0
		}
	}

	for end_i := 0; i < len(si.str); i++ {
		if si.str[i] == si.stop[end_i] {
			end_i++
			if end_i == len(si.stop) {
				endIdx = i - len(si.stop) + 1
				break
			}
		} else {
			end_i = 0
		}
	}

	si.i = i
	if i == len(si.str) && endIdx == -1 {
		si.prefix_end = i
		return "", false
	}

	return si.str[startIdx:endIdx], true
}

func (si *stringIterator) Prefix() string {
	return si.str[si.prefix_begin:si.prefix_end]
}

func (si *stringIterator) GetWholePrefix() string {
	return si.str[:si.prefix_end]
}

func formatError(format, accessor string) error {
	if accessor == "http-body" {
		format += ". Did you forget (fetch: json)?"
	}
	return fmt.Errorf(format, accessor)
}

func access(accessor string, storage stdlib_types.Storage) (string, error) {
	//TODO: nested
	iter := stringIterator{str: accessor, start: "[", stop: "]"}
	accessor, ok := iter.Next()
	var current any
	current = storage[iter.Prefix()]

	for ; ok; accessor, ok = iter.Next() {
		dict, isMap := current.(map[string]any)
		if isMap {
			current = dict[accessor]
			if current == nil {
				return "", formatError("value with key: (%s) does not exist", accessor)
			}
		} else {
			array, isArray := current.([]any)
			if !isArray {
				return "", formatError("failed to dereference by key (%s), previous values was neither a map[string]any nor an []any", accessor)
			}
			idx, err := strconv.ParseInt(accessor, 0, 64)
			if err != nil {
				return "", formatError("failed to access an array! expected int as a key but got: (%s)", accessor)
			}

			if int(idx) >= len(array) {
				//TODO: what array? print whole prefix
				return "", fmt.Errorf("out of bounds! Index: %d. Len: %d", idx, len(array))
			}

			current = array[idx]
		}
	}

	res, ok := current.(string)
	if !ok {
		return "", fmt.Errorf("injection result expected to be a string but it is of type (%T)", current)
	}
	return res, nil
}

func TryInject(injectee string, storage stdlib_types.Storage) (string, error) {
	iter := stringIterator{str: injectee, start: "${", stop: "}"}
	result := strings.Builder{}

	accessor, ok := iter.Next()

	for ; ok; accessor, ok = iter.Next() {
		result.WriteString(iter.Prefix())
		val, err := access(accessor, storage)
		if err != nil {
			return "", err
		}
		result.WriteString(val)
	}
	result.WriteString(iter.Prefix())
	return result.String(), nil
}
