package stdlib_types

type Storage map[string]map[string]any

type Verdict string

const (
	NT Verdict = "NOT TESTED"
	OK Verdict = "ACCEPTED"
	WA Verdict = "WRONG ANSWER"
	PE Verdict = "PRESENTATION ERROR"
	EF Verdict = "EPIC FAIL"
	CF Verdict = "CHECK FAILED"
	TL Verdict = "TIME LIMIT EXCEED"
	ML Verdict = "MEMORY LIMIT EXCEED"
	SV Verdict = "SECURITY VIOLATION"
	CE Verdict = "COMPILATION ERROR"
	RE Verdict = "RUNTIME ERROR"
	IO Verdict = "INVALID IO"
	TT Verdict = "TESTED"
	WT Verdict = "WALL TIME LIMIT"
	DE Verdict = "DOWNLOAD ERROR"
	MC Verdict = "MISCONFIGURED"
)

type ExecutorHandler func(
	id string,
	args any,
	storage Storage,
) ExecutorResult

func (s Storage) InitSection(sectionName string) {
	s[sectionName] = make(map[string]any)
}

type ExecutorResult struct {
	Verdict Verdict
	Comment string
}

func (er *ExecutorResult) AppendSectionName(section string) {
	er.Comment = "In section (" + section + "): " + er.Comment
}
