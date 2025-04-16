package executors

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	shared "backend-testing-module-shared"
)

func HandleCookiesExecutor(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	handleCookies, ok := args.(*bool)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "HANDLE_COOKIES executor accepts single bool argument",
		}
	}

	stdlib_helpers.DefaultHttpClient.ToogleCookiesHandling(*handleCookies)

	return stdlib_types.ExecutorResult{
		Verdict: shared.OK,
	}
}
