package checkercontext

import (
	stdlib_types "backend-testing-module-checker/stdlib/types"
	shared "backend-testing-module-shared"
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type executableSection struct {
	Name      string
	Executors ExecutorCollection
	Checkers  ExecutorCollection
}

type Test struct {
	Requests  SectionCollection
	Responses SectionCollection
}

type checkerResult struct {
	stdlib_types.ExecutorResult
	Log OutputLogStructure
}

func (t *Test) Execute() {
	sections, err := t.newExecutableSections()
	if err != nil {
		writeAnswer(stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: err.Error(),
		})
		return
	}

	logger := Logger()
	for _, section := range sections {
		logger.StartNewSection(section.Name)
		context.storage.InitSection(section.Name)
		for _, executor := range section.Executors {
			result := context.CallExecutor(section.Name, executor.Name, executor.Args)
			if result.Verdict != shared.OK {
				result.AppendSectionName(section.Name)
				writeAnswer(result)
				return
			}
		}

		for _, checker := range section.Checkers {
			result := context.CallChecker(section.Name, checker.Name, checker.Args)
			if result.Verdict != shared.OK {
				result.AppendSectionName(section.Name)
				writeAnswer(result)
				return
			}
		}
	}

	writeAnswer(stdlib_types.ExecutorResult{
		Verdict: shared.OK,
		Comment: "OK",
	})
}

func writeAnswer(result stdlib_types.ExecutorResult) {
	checkerResult := checkerResult{
		ExecutorResult: result,
		Log:            *Logger().getOutput(),
	}

	json.NewEncoder(os.Stdout).Encode(checkerResult)
}

func ParseTest(testPath string, asnwerPath string) (*Test, error) {
	test := &Test{}
	data, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open testcase file: %v", err)
	}
	err = yaml.Unmarshal(data, &test.Requests)
	if err != nil {
		return nil, fmt.Errorf("failed to parse testcase: %v", err)
	}

	data, err = os.ReadFile(asnwerPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open answer file: %v", err)
	}
	err = yaml.Unmarshal(data, &test.Responses)
	if err != nil {
		return nil, fmt.Errorf("failed to parse answer: %v", err)
	}

	return test, nil
}

func (t *Test) newExecutableSections() ([]executableSection, error) {
	var (
		sectionMap = make(map[string]*executableSection)
		es         []executableSection
	)

	for _, section := range t.Requests {
		es = append(es, executableSection{Name: section.Id, Executors: section.Executors})
		sectionMap[section.Id] = &es[len(es)-1]
	}

	for _, section := range t.Responses {
		execSection, ok := sectionMap[section.Id]
		if !ok {
			return nil, fmt.Errorf("asnwer file contains section %s which does not exist in test file", section.Id)
		}

		execSection.Checkers = section.Executors
	}

	return es, nil

}

func WriteError(msg string) {
	writeAnswer(stdlib_types.ExecutorResult{
		Verdict: shared.CF,
		Comment: msg,
	})
}
