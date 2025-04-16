package executors

import (
	stdlib_helpers "backend-testing-module-checker/stdlib/helpers"
	stdlib_types "backend-testing-module-checker/stdlib/types"
	shared "backend-testing-module-shared"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

const (
	STATUS_CODE_KEY  = "http-status-code"
	HEADERS_KEY      = "http-headers"
	BODY_KEY         = "http-body"
	RAW_RESPONSE_KEY = "http-raw-response"
)

type HttpArgs struct {
	Endpoint      string              `yaml:"endpoint"`
	Body          string              `yaml:"body"`
	CustomHeaders map[string][]string `yaml:"custom-headers"`
	Method        string              `yaml:"method"`
	FetchFormat   string              `yaml:"fetch"`
	QueryParams   map[string][]string `yaml:"query-params"`
}

func (ha *HttpArgs) Inject(storage stdlib_types.Storage) error {
	var err error

	ha.Endpoint, err = stdlib_helpers.TryInject(ha.Endpoint, storage)
	if err != nil {
		return fmt.Errorf("injection into 'endpoint' failed: %v", err)
	}

	ha.Body, err = stdlib_helpers.TryInject(ha.Body, storage)
	if err != nil {
		return fmt.Errorf("injection into 'body' failed: %v", err)
	}

	ha.Method, err = stdlib_helpers.TryInject(ha.Method, storage)
	if err != nil {
		return fmt.Errorf("injection into 'method' failed: %v", err)
	}

	for key, list := range ha.CustomHeaders {
		for i, val := range list {
			list[i], err = stdlib_helpers.TryInject(val, storage)
			if err != nil {
				return fmt.Errorf("injection into 'custom-header[%s]' failed: %v", key, err)
			}
		}
	}

	for key, list := range ha.QueryParams {
		for i, val := range list {
			list[i], err = stdlib_helpers.TryInject(val, storage)
			if err != nil {
				return fmt.Errorf("injection into 'query-params[%s]' failed: %v", key, err)
			}
		}
	}

	return nil
}

func HttpExecutor(
	id string,
	args any,
	storage stdlib_types.Storage,
) stdlib_types.ExecutorResult {
	arguments, ok := args.(*HttpArgs)
	if !ok {
		return stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: "failed to retrieve arguments",
		}
	}

	err := arguments.Inject(storage)
	if err != nil {
		return stdlib_types.ExecutorResult{
			Verdict: shared.PE,
			Comment: err.Error(),
		}
	}

	resp, err := stdlib_helpers.DefaultHttpClient.NewRequest().
		SetMethod(arguments.Method).
		SetEndpoint(arguments.Endpoint).
		SetBody(arguments.Body).
		SetCustomHeaders(arguments.CustomHeaders).
		SetQueryParams(arguments.QueryParams).
		Do()

	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return stdlib_types.ExecutorResult{
					Verdict: shared.TL,
					Comment: fmt.Sprintf("request timed out: %v", err),
				}
			}

			if opError, ok := urlErr.Unwrap().(*net.OpError); ok {
				return stdlib_types.ExecutorResult{
					Verdict: shared.MC,
					Comment: fmt.Sprintf("operation error: %v", opError),
				}
			}
		}
		return stdlib_types.ExecutorResult{
			Verdict: shared.PE,
			Comment: fmt.Sprintf("error during request: %v", err),
		}
	}

	storage[id][RAW_RESPONSE_KEY] = resp
	storage[id][STATUS_CODE_KEY] = resp.StatusCode
	storage[id][HEADERS_KEY] = resp.Header
	return *bodyDecode(arguments.FetchFormat, storage[id], resp)

}

func bodyDecode(format string, sectionStorage map[string]any, resp *http.Response) *stdlib_types.ExecutorResult {
	if format == "" {
		return &stdlib_types.ExecutorResult{Verdict: shared.OK}
	}
	var bodyObject any
	decoder, err := stdlib_helpers.NewDecodeMatcher(format)
	if err != nil {
		return &stdlib_types.ExecutorResult{
			Verdict: shared.EF,
			Comment: fmt.Sprintf("failed to create decoder: %v", err),
		}
	}

	err = decoder.Decode(resp.Body, &bodyObject)
	if err != nil {
		return &stdlib_types.ExecutorResult{
			Verdict: shared.PE,
			Comment: fmt.Sprintf("couldn't decode response body: %v. Response status: %v", err, resp.Status),
		}
	}

	sectionStorage[BODY_KEY] = bodyObject
	return &stdlib_types.ExecutorResult{Verdict: shared.OK}
}
