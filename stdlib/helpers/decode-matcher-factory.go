package stdlib_helpers

import (
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"fmt"
)

func NewDecodeMatcher(format string) (stdlib_types.DecodeMather, error) {
	switch format {
	case "json", "":
		return new(JsonConverter), nil

	default:
		return nil, fmt.Errorf("%s is not supported format", format)
	}
}
