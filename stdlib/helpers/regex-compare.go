package stdlib_helpers

import (
	"fmt"
	"regexp"
	"strings"
)

func IsRegex(str string) bool {
	return strings.HasPrefix(str, "$re{") && strings.HasSuffix(str, "}")
}

func ExtractRegex(str string) (*regexp.Regexp, error) {

	if str == "" {
		return nil, fmt.Errorf("empty string is not a regex")
	}

	regex, found := strings.CutPrefix(str, "$re{")
	if !found {
		return nil, fmt.Errorf("%v is not a regex", str)
	}

	regex, found = strings.CutSuffix(regex, "}")
	if !found {
		return nil, fmt.Errorf("%v is not a regex", regex)
	}

	return regexp.Compile("^" + regex + "$")
}

func TryRegexCompare(a, b string) (bool, error) {
	var (
		regex         *regexp.Regexp
		stringToMatch string
		err           error
	)

	if IsRegex(a) {
		stringToMatch = b
		regex, err = ExtractRegex(a)
		if err != nil {
			return false, err
		}
	}

	if IsRegex(b) {
		if regex != nil {
			return false, fmt.Errorf("both strings can not be regexes")
		}
		stringToMatch = a
		regex, err = ExtractRegex(b)
		if err != nil {
			return false, err
		}
	}

	if regex == nil {
		return a == b, nil
	}

	return regex.Match([]byte(stringToMatch)), nil

}
