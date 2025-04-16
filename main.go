package main

import (
	checkercontext "backend-testing-module-checker/checker-context"
	"backend-testing-module-checker/cmdline"
	executorexample "backend-testing-module-checker/executor-example"
	"backend-testing-module-checker/stdlib"
	"backend-testing-module-checker/stdlib/checkers"
	"backend-testing-module-checker/stdlib/executors"
	"fmt"
)

func main() {
	args, err := cmdline.ParseCmdlineArgs()
	if err != nil {
		checkercontext.WriteError(
			fmt.Sprintf("failed to parse cmdline arguments: %v", err),
		)
		return
	}
	stdlib.InitStdlib()

	checkercontext.AddChecker[executorexample.LessThanArgs]("LESS_THAN", executorexample.LessThanHandler)
	checkercontext.AddExecutor[executors.DBQueryArgs]("DB_QUERY", executors.DBQueryHandler)
	checkercontext.AddChecker[checkers.ValuesMatchArgs]("MATCH_VALUES", checkers.ValuesMatchChecker)

	testcase, err := checkercontext.ParseTest(args.TestFile, args.AnswerFile)
	if err != nil {
		checkercontext.WriteError(err.Error())
		return
	}
	testcase.Execute()
}
