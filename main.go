package main

import (
	checkercontext "backend-testing-module-checker/checker-context"
	"backend-testing-module-checker/cmdline"
	"backend-testing-module-checker/stdlib"
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

	testcase, err := checkercontext.ParseTest(args.TestFile, args.AnswerFile)
	if err != nil {
		checkercontext.WriteError(err.Error())
		return
	}
	testcase.Execute()
}
