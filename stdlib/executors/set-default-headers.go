package executors

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
)

func SetDefaultHeadersExecutor(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	headers, ok := args.(*map[string][]string)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "Can't cast SET_DEFAULT_HEADERS arguments",
		}
	}

	stdlib_helpers.DefaultHttpClient.AppendHeaders(*headers)

	return stdlib_types.ExecutorResult{
		Verdict: stdlib_types.OK,
	}
}
