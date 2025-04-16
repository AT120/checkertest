package stdlib_helpers

import (
	"io"
	"net/http"
	"strings"
)

type RequestBuilder struct {
	client *HttpClient

	method   string
	endpoint string
	body     io.Reader
	headers  map[string][]string
	query    map[string][]string
}

func (builder *RequestBuilder) SetMethod(method string) *RequestBuilder {
	builder.method = method
	return builder
}

func (builder *RequestBuilder) SetEndpoint(endpoint string) *RequestBuilder {
	builder.endpoint = endpoint
	return builder
}

func (builder *RequestBuilder) SetBody(body string) *RequestBuilder {
	builder.body = strings.NewReader(body)
	return builder
}

func (builder *RequestBuilder) SetCustomHeaders(headers map[string][]string) *RequestBuilder {
	builder.headers = headers
	return builder
}

func (builder *RequestBuilder) SetQueryParams(query map[string][]string) *RequestBuilder {
	builder.query = query
	return builder
}

func (builder *RequestBuilder) Do() (*http.Response, error) { //TODO: rename
	return builder.client.makeBuildedRequest(builder)
}
