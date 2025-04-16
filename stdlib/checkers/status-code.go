package checkers

import (
	"backend-testing-module-checker/stdlib/executors"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"fmt"
)

type StatusCodeCheckerArgs struct {
	Exact     int `yaml:"exact"`
	RangeFrom int `yaml:"range-from"`
	RangeTo   int `yaml:"range-to"`
}

func (args *StatusCodeCheckerArgs) Validate() error {
	if (args.Exact != 0) && (args.RangeTo != 0) {
		return fmt.Errorf("'exact' and 'range' arguments should not be combined")
	}

	if (args.Exact == 0) && (args.RangeTo == 0) {
		return fmt.Errorf("either 'exact' or 'range' arguments have to be set")
	}

	return nil
}

func (args *StatusCodeCheckerArgs) WAMessage(statusCode int) string {
	if args.Exact != 0 {
		return fmt.Sprintf("Wrong status code! Got: %d; Wanted: %d", statusCode, args.Exact)
	} else {
		return fmt.Sprintf("Wrong status code! Got: %d; Wanted between %d and %d", statusCode, args.RangeFrom, args.RangeTo)
	}

}

func StatusCodeChecker(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	arguments, ok := args.(*StatusCodeCheckerArgs)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "HTTP executor failed to retrieve arguments",
		}
	}

	err := arguments.Validate()
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: err.Error(),
		}
	}

	statusCodeUntyped, ok := storage[id][executors.STATUS_CODE_KEY]
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "status code was not set for this section",
		}
	}
	statusCode, ok := statusCodeUntyped.(int)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "status code was written to storage with invalid type",
		}
	}

	if statusCode == arguments.Exact || (statusCode >= arguments.RangeFrom && statusCode < arguments.RangeTo) {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.OK,
			Comment: "OK",
		}
	}

	return stdlib_types.ExecutorResult{
		Verdict: stdlib_types.WA,
		Comment: arguments.WAMessage(statusCode),
	}
}
