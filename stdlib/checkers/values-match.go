package checkers

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"fmt"
	"strings"
)

type ValuesMatchArgs struct {
	Format      string `yaml:"format"`
	Strict      bool   `yaml:"strict"`
	FirstValue  string `yaml:"first_value"`
	SecondValue string `yaml:"second_value"`
}

func tryInjectValues(args *ValuesMatchArgs, storage stdlib_types.Storage) error {
	var err error
	args.Format, err = stdlib_helpers.TryInject(args.Format, storage)
	if err != nil {
		return fmt.Errorf(`"format" field injection failed: %v`, err)
	}

	args.FirstValue, err = stdlib_helpers.TryInject(args.FirstValue, storage)
	if err != nil {
		return fmt.Errorf(`"first_string" field injection failed: %v`, err)
	}

	args.SecondValue, err = stdlib_helpers.TryInject(args.SecondValue, storage)
	if err != nil {
		return fmt.Errorf(`"second_string" field injection failed: %v`, err)
	}

	return nil
}

func ValuesMatchChecker(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	arguments, ok := args.(*ValuesMatchArgs)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "MATCH_VALUES checker failed to retrieve arguments",
		}
	}

	err := tryInjectValues(arguments, storage)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: "MATCH_VALUES: " + err.Error(),
		}
	}

	decodeMatcher, err := stdlib_helpers.NewDecodeMatcher(arguments.Format)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: err.Error(),
		}
	}

	var firstObject, secondObject any
	err = decodeMatcher.Decode(strings.NewReader(arguments.FirstValue), &firstObject)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "MATCH_VALUES checker failed to decode first string: " + err.Error(),
		}
	}

	err = decodeMatcher.Decode(strings.NewReader(arguments.SecondValue), &secondObject)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "MATCH_VALUES checker failed to decode second string: " + err.Error(),
		}
	}

	var result bool
	if arguments.Strict {
		result, err = decodeMatcher.StrictlyCompare(firstObject, secondObject)
	} else {
		result, err = decodeMatcher.LooselyCompare(firstObject, secondObject)
	}

	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "MATCH_VALUES: match failed: " + err.Error(),
		}
	}

	if result {
		return stdlib_types.ExecutorResult{Verdict: stdlib_types.OK}
	}

	return stdlib_types.ExecutorResult{
		Verdict: stdlib_types.WA,
		Comment: "Values do not match; FirstValue: " + fmt.Sprintf("%v", firstObject) + "; SecondValue: " + fmt.Sprintf("%v", secondObject),
	}
}
