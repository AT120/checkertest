package executors

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	shared "backend-testing-module-shared"
)

func SetBaseUrlExecutor(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	baseUrl, ok := args.(*string)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "Couldn't cast SET_BASE_URL arguments",
		}
	}

	stdlib_helpers.DefaultHttpClient.SetBaseUrl(*baseUrl)

	return stdlib_types.ExecutorResult{
		Verdict: shared.OK,
	}
}
