package stdlib_types

import shared "backend-testing-module-shared"

type Storage map[string]map[string]any

type ExecutorHandler func(
	id string,
	args any,
	storage Storage,
) ExecutorResult

func (s Storage) InitSection(sectionName string) {
	s[sectionName] = make(map[string]any)
}

type ExecutorResult struct {
	Verdict shared.Verdict
	Comment string
}

func (er *ExecutorResult) AppendSectionName(section string) {
	er.Comment = "In section (" + section + "): " + er.Comment
}
