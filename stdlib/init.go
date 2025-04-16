package stdlib

import (
	checkercontext "backend-testing-module-checker/checker-context"
	"backend-testing-module-checker/stdlib/checkers"
	"backend-testing-module-checker/stdlib/executors"
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	"net/http"
)

func InitStdlib() {
	stdlib_helpers.DefaultHttpClient.InitHttp(http.DefaultClient)

	checkercontext.AddExecutor[executors.HttpArgs]("HTTP", executors.HttpExecutor)
	checkercontext.AddExecutor[bool]("HANDLE_COOKIES", executors.HandleCookiesExecutor)
	checkercontext.AddExecutor[map[string][]string]("SET_DEFAULT_HEADERS", executors.SetDefaultHeadersExecutor)
	checkercontext.AddExecutor[string]("SET_BASE_URL", executors.SetBaseUrlExecutor)

	checkercontext.AddChecker[checkers.StatusCodeCheckerArgs]("STATUS_CODE", checkers.StatusCodeChecker)
	checkercontext.AddChecker[checkers.BodyMatchArgs]("BODY_MATCH", checkers.BodyMatchChecker)
}
