package executorexample

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	"fmt"
	"strconv"
)

type LessThanArgs struct {
	A string `yaml:"a"`
	B string `yaml:"b"`
}

// сигнатура соответсвует stdlib_types.ExecutorHandler
func LessThanHandler(id string, args any, storage stdlib_types.Storage) stdlib_types.ExecutorResult {

	// Кастим переданные аргументы в нашу структуру
	arguments, ok := args.(*LessThanArgs)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.EF,
			Comment: "failed to retrieve arguments. Checker bug",
		}
	}

	// выполняем подстановку (разворачиваем конструкции вида ${section[smth][smth]})
	aStr, err := stdlib_helpers.TryInject(arguments.A, storage)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("could not parse A argument: %v", err),
		}
	}

	// парсим число
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: "Failed to parse float",
		}
	}

	// выполняем подстановку для B
	bStr, err := stdlib_helpers.TryInject(arguments.B, storage)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: fmt.Sprintf("could not parse B argument: %v", err),
		}
	}

	// парсим B
	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: stdlib_types.PE,
			Comment: "Failed to parse float",
		}
	}

	// проверяем условие
	if a < b {
		return stdlib_types.ExecutorResult{Verdict: stdlib_types.OK}
	}

	// иначе возвращаем ошибку
	return stdlib_types.ExecutorResult{
		Verdict: stdlib_types.WA,
		Comment: fmt.Sprintf("%f >= %f, expected otherwise", a, b),
	}
}
