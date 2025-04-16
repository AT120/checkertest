package stdlib_helpers

import (
	checkercontext "backend-testing-module-checker/checker-context"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type HttpClient struct {
	client     *http.Client
	headers    http.Header
	rawBaseUrl string
	cookieJar  *cookiejar.Jar
}

var DefaultHttpClient = HttpClient{}

type RequestResponse struct {
	request *http.Request
}

func (c *HttpClient) InitHttp(newClient *http.Client) {
	if newClient != nil {
		c.client = newClient
	} else {
		c.client = &http.Client{}
	}

	if c.client.Jar == nil {
		c.cookieJar, _ = cookiejar.New(nil)
		c.client.Jar = c.cookieJar
	}

	if c.client.Transport != nil {
		c.client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		logResponse(req.Response)
		logRequest(req)

		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}

	c.headers = make(http.Header)
}

func (c *HttpClient) SetBaseUrl(rawURL string) error {
	if !strings.HasPrefix(rawURL, "http") {
		return fmt.Errorf("base url has to start with \"http\"")
	}
	c.rawBaseUrl = rawURL

	return nil
}

func (c *HttpClient) AppendHeaders(headers map[string][]string) {
	for key, valueList := range headers {
		for _, value := range valueList {
			c.headers.Add(key, value)
		}
	}
}

func (c *HttpClient) NewRequest() *RequestBuilder {
	return &RequestBuilder{client: c}
}

func (c *HttpClient) ToogleCookiesHandling(handle bool) {
	logger := checkercontext.Logger()
	if handle {
		logger.Printf("! CHECKER ENABLED COOKIES HANDLING")
		c.client.Jar = c.cookieJar
	} else {
		logger.Printf("! CHECKER DISABLED COOKIES HANDLING")
		c.client.Jar = nil
	}
}

func (c *HttpClient) makeBuildedRequest(builder *RequestBuilder) (*http.Response, error) {
	var err error
	reqUrl := builder.endpoint
	if c.rawBaseUrl != "" {
		reqUrl, err = url.JoinPath(c.rawBaseUrl, builder.endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to join base url with endpoint: %v", err)
		}
	}

	req, err := http.NewRequest(builder.method, reqUrl, builder.body)
	if err != nil {
		return nil, err
	}

	for key, list := range builder.headers {
		for _, value := range list {
			req.Header.Add(key, value)
		}
	}

	for key, list := range c.headers {
		for _, value := range list {
			req.Header.Add(key, value)
		}
	}

	queries := req.URL.Query()
	for key, list := range builder.query {
		for _, value := range list {
			queries.Add(key, value)
		}
	}

	req.URL.RawQuery = queries.Encode()
	logRequest(req)

	resp, err := c.client.Do(req)

	logResponse(resp)

	return resp, err
}

func beautifyBody(body []byte) (string, bool) {
	var beautifiedBody bytes.Buffer
	err := json.Indent(&beautifiedBody, body, "", "\t")
	if err != nil {
		return string(body), false
	}
	return beautifiedBody.String(), true
}

func logRequest(req *http.Request) {
	logger := checkercontext.Logger()

	var headers []string
	for key, value := range req.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", key, strings.Join(value, "; ")))
	}

	var bodyData []byte
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err == nil {
			bodyData, _ = io.ReadAll(body)
		}
	}

	logger.WriteNewHttpRequest(&checkercontext.HttpRequest{
		HttpVersion: req.Proto,
		Method:      req.Method,
		URL:         req.URL.String(),
		Headers:     headers,
		Body:        string(bodyData),
	})
}

func logResponse(resp *http.Response) {
	logger := checkercontext.Logger()
	if resp == nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		body = []byte("!!!FAILED TO READ REQUEST BODY!!!")
	}

	var headers []string
	for key, value := range resp.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", key, strings.Join(value, "; ")))
	}

	logger.WriteHttpResponse(&checkercontext.HttpResponse{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       string(body),
	})

}
