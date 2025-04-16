package checkers

import (
	"backend-testing-module-checker/stdlib/executors"
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	shared "backend-testing-module-shared"
	"fmt"
	"strings"
)

type BodyMatchArgs struct {
	Format  string `yaml:"format"`
	Strict  bool   `yaml:"strict"`
	Pattern string `yaml:"pattern"`
}

func tryInject(args *BodyMatchArgs, storage stdlib_types.Storage) error {
	var err error
	args.Format, err = stdlib_helpers.TryInject(args.Format, storage)
	if err != nil {
		return fmt.Errorf(`"format" field injection failed: %v`, err)
	}

	args.Pattern, err = stdlib_helpers.TryInject(args.Pattern, storage)
	if err != nil {
		return fmt.Errorf(`"pattern" field injection failed: %v`, err)
	}

	return nil
}

func BodyMatchChecker(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	arguments, ok := args.(*BodyMatchArgs)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "BODY_MATCH checker failed to retrieve arguments",
		}
	}

	err := tryInject(arguments, storage)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: shared.PE,
			Comment: "BODY_MATCH: " + err.Error(),
		}
	}

	requestBody, ok := storage[id][executors.BODY_KEY]
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "BODY_MATCH checker failed to get body of previous request. " +
				"Make sure there was an HTTP executor in a test file in the same section",
		}
	}

	decodeMatcher, err := stdlib_helpers.NewDecodeMatcher(arguments.Format)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: err.Error(),
		}
	}

	var patternObject any
	err = decodeMatcher.Decode(strings.NewReader(arguments.Pattern), &patternObject)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "BODY_MATCH checker failed to decode pattern: " + err.Error(),
		}
	}

	var result bool
	if arguments.Strict {
		result, err = decodeMatcher.StrictlyCompare(requestBody, patternObject)
	} else {
		result, err = decodeMatcher.LooselyCompare(requestBody, patternObject)
	}

	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "BODY_MATCH: match failed: " + err.Error(),
		}
	}

	if result {
		return stdlib_types.ExecutorResult{Verdict: shared.OK}
	}

	return stdlib_types.ExecutorResult{
		Verdict: shared.WA,
		Comment: "Response body does not match expected pattern",
	}

}
