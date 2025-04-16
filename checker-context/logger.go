package checkercontext

import (
	"encoding/json"
	"fmt"
)

type CheckerLogger struct {
	output              OutputLogStructure
	lastHttpOperationId int
}

var logger *CheckerLogger = &CheckerLogger{}

const separator = "-------"

func Logger() *CheckerLogger {
	return logger

}

func (cl *CheckerLogger) StartNewSection(sectionName string) {
	cl.output.Sections = append(cl.output.Sections, Section{
		Name: sectionName,
	})
}

func (cl *CheckerLogger) lastSection() *Section {
	return &cl.output.Sections[len(cl.output.Sections)-1]
}

func (cl *CheckerLogger) WriteNewHttpRequest(request *HttpRequest) {
	section := cl.lastSection()
	section.Operations = append(section.Operations, Operation{
		Type:    HTTP,
		Details: &HttpOperationDetails{},
	})
	op := section.lastOperation()
	details, ok := op.Details.(*HttpOperationDetails)
	if !ok {
		return
	}

	cl.lastHttpOperationId = len(section.Operations) - 1
	details.Request = *request
}

func (cl *CheckerLogger) WriteHttpResponse(response *HttpResponse) {
	section := cl.lastSection()
	if len(section.Operations) < cl.lastHttpOperationId {
		cl.Printf("CHECKER BUG: http request-response mismatch")
		return
	}

	details, ok := section.Operations[cl.lastHttpOperationId].Details.(*HttpOperationDetails)
	if !ok {
		cl.Printf("CHECKER BUG: invalid operation id for http response")
		return
	}

	details.Response = *response
}

func (cl *CheckerLogger) Printf(format string, args ...any) {
	section := cl.lastSection()
	section.Operations = append(section.Operations, Operation{
		Type:    LOG,
		Details: fmt.Sprintf(format, args...),
	})
}

func (cl *CheckerLogger) getOutput() *OutputLogStructure {
	return &cl.output
}

func (cl *CheckerLogger) string() string {
	output, _ := json.Marshal(&cl.output)
	return string(output)
}
