package executors

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
)

func HandleCookiesExecutor(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	handleCookies, ok := args.(*bool)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "HANDLE_COOKIES executor accepts single bool argument",
		}
	}

	stdlib_helpers.DefaultHttpClient.ToogleCookiesHandling(*handleCookies)

	return stdlib_types.ExecutorResult{
		Verdict: stdlib_types.OK,
	}
}
