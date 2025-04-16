package checkercontext

type OperationType string

const (
	HTTP OperationType = "HTTP"
	LOG  OperationType = "LOG"
)

type OutputLogStructure struct {
	Sections []Section
}

type Section struct {
	Name       string
	Operations []Operation
}

type Operation struct {
	Type    OperationType
	Details any
}

type HttpOperationDetails struct {
	Request  HttpRequest
	Response HttpResponse
}

type HttpResponse struct {
	Headers    []string
	Body       string
	StatusCode int
}

type HttpRequest struct {
	HttpVersion string
	Method      string
	URL         string
	Headers     []string
	Body        string
}

type LogOperationDetails string

func (section *Section) lastOperation() *Operation {
	return &section.Operations[len(section.Operations)-1]
}
