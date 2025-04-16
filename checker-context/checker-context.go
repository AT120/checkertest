package checkercontext

import (
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"fmt"
)

type checkerContext struct {
	RequestExecutors map[string]stdlib_types.ExecutorHandler
	ResponseCheckers map[string]stdlib_types.ExecutorHandler

	storage     stdlib_types.Storage
	argsFactory map[string]func() any
}

var context = checkerContext{
	argsFactory:      make(map[string]func() any),
	RequestExecutors: make(map[string]stdlib_types.ExecutorHandler),
	ResponseCheckers: make(map[string]stdlib_types.ExecutorHandler),
	storage:          make(stdlib_types.Storage),
}

func AddExecutor[ArgType any](name string, handler stdlib_types.ExecutorHandler) {
	context.argsFactory[name] = func() any {
		var arg ArgType
		return &arg
	}
	// context.argsFactory[name] =
	context.RequestExecutors[name] = handler
}

func AddChecker[ArgType any](name string, handler stdlib_types.ExecutorHandler) {
	context.argsFactory[name] = func() any {
		var arg ArgType
		return &arg
	}
	// context.argsFactory[name] =
	context.ResponseCheckers[name] = handler
}

func (c *checkerContext) CallExecutor(id string, name string, args any) stdlib_types.ExecutorResult {
	executor := c.RequestExecutors[name]
	if executor == nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: fmt.Sprintf("executor with name %s does not exist", name),
		}
	}

	return executor(id, args, c.storage)
}

func (c *checkerContext) CallChecker(id string, name string, args any) stdlib_types.ExecutorResult {
	checker := c.ResponseCheckers[name]
	if checker == nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: fmt.Sprintf("checker with name %s does not exist", name),
		}
	}

	return checker(id, args, c.storage)
}
