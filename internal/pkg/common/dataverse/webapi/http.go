package webapi

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPRequestError is used for errors found during composing an HTTP request.
type HTTPRequestError struct {
	err error
}

// HTTPRequest - T is the type of the decoded response body.
type HTTPRequest[T any] struct {
	result         *T
	client         *Client
	method         string
	path           string
	query          url.Values
	header         http.Header
	body           any
	expectedStatus int
	onSuccess      []func(ctx context.Context, result *T) error
}

type NoResult struct{}

type RawResponse struct {
	Response *http.Response
	Body     io.Reader
}

type httpRequest interface {
	Do(ctx context.Context) error
	prepareRequest(ctx context.Context) (*http.Request, error)
	processResponse(ctx context.Context, resp *http.Response) error
}

func NewHTTPRequest[T any](result *T, client *Client, method, path string, body any) *HTTPRequest[T] {
	return &HTTPRequest[T]{
		result:         result,
		client:         client,
		method:         method,
		path:           path,
		query:          make(url.Values),
		header:         make(http.Header),
		body:           body,
		expectedStatus: http.StatusOK,
	}
}

func NewHTTPRequestError(err error) *HTTPRequestError {
	return &HTTPRequestError{err: err}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.config.DebugRequest {
		dumpRequest(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.config.DebugResponse {
		dumpResponse(req, resp)
	}

	return resp, nil
}

func (e *HTTPRequestError) Do(_ context.Context) error {
	return e.err
}

func (e *HTTPRequestError) prepareRequest(context.Context) (*http.Request, error) {
	return nil, e.err
}

func (*HTTPRequestError) processResponse(context.Context, *http.Response) error {
	panic(errors.New("shouldn't be called"))
}

func (r *HTTPRequest[T]) OnSuccess(fn func(ctx context.Context, result *T) error) *HTTPRequest[T] {
	r.onSuccess = append(r.onSuccess, fn)
	return r
}

func (r *HTTPRequest[T]) Header(key, value string) *HTTPRequest[T] {
	r.header.Set(key, value)
	return r
}

func (r *HTTPRequest[T]) Select(keys ...string) *HTTPRequest[T] {
	r.query.Set("$select", strings.Join(keys, ","))
	return r
}

func (r *HTTPRequest[T]) Filter(filter string) *HTTPRequest[T] {
	r.query.Set("$filter", filter)
	return r
}

func (r *HTTPRequest[T]) ExpectStatus(status int) *HTTPRequest[T] {
	r.expectedStatus = status
	return r
}

func (r *HTTPRequest[T]) Do(ctx context.Context) error {
	req, err := r.prepareRequest(ctx)
	if err != nil {
		return err
	}

	// Perform the HTTP request
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("dataverse request failed: %w", err)
	}

	// Defer closing the response body to prevent resource leaks.
	// The close error if there is no other error.
	defer func() {
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body) // drain body to reuse connection
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}
	}()

	return r.processResponse(ctx, resp)
}

func (r *HTTPRequest[T]) prepareRequest(ctx context.Context) (*http.Request, error) {
	// Compose URL
	absURL := r.client.baseURL.ResolveReference(&url.URL{Path: r.path, RawQuery: r.query.Encode()})

	// Body
	body, contentType, err := r.bodyReader()
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, r.method, absURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("cannot create dataverse request: %w", err)
	}

	// Set headers
	req.Header = r.header.Clone()

	req.Header.Set("Accept", ContentTypeJSON)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

// processResponse validates the HTTP response and decodes the JSON body into the result type.
func (r *HTTPRequest[T]) processResponse(ctx context.Context, resp *http.Response) (err error) {
	// In batch/change set requests, content is not returned, so the Location header must be followed.
	if r.shouldFollowLocationHeader(resp) {
		return r.followLocationHeader(ctx, resp)
	}

	// Return if there is no content
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// Decompress body if needed
	reader, err := decompressBody(resp)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := reader.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	// Check status code
	if resp.StatusCode != r.expectedStatus {
		errResult := APIError{}
		_ = decodeBody(resp, reader, &errResult)
		return UnexpectedStatusError{
			Method:   resp.Request.Method,
			URL:      resp.Request.URL.String(),
			Expected: r.expectedStatus,
			Actual:   resp.StatusCode,
			Message:  errResult.Error.Message,
		}
	}

	// Decode body
	if err = decodeBody(resp, reader, r.result); err != nil {
		return err
	}

	// Call success callbacks
	for _, fn := range r.onSuccess {
		if err := fn(ctx, r.result); err != nil {
			return err
		}
	}

	// Return a pointer to the decoded result.
	return nil
}

func (r *HTTPRequest[T]) bodyReader() (_ io.Reader, contentType string, _ error) {
	switch body := r.body.(type) {
	case nil:
		return nil, "", nil
	case io.ReadSeeker:
		return body, "", nil
	case io.Reader:
		panic(errors.New("expected io.ReadSeeker, found io.Reader: io.ReadSeeker is required to rewind body before retry"))
	default:
		bytes, err := json.Marshal(r.body)
		if err != nil {
			return nil, "", fmt.Errorf("cannot marshal body: %w", err)
		}
		return io.NopCloser(strings.NewReader(string(bytes))), ContentTypeJSON, nil
	}
}

func (r *HTTPRequest[T]) shouldFollowLocationHeader(resp *http.Response) bool {
	location := resp.Header.Get("Location")
	return location != "" && resp.StatusCode == http.StatusNoContent && r.expectedStatus != http.StatusNoContent
}

func (r *HTTPRequest[T]) followLocationHeader(ctx context.Context, resp *http.Response) error {
	locationURL, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return fmt.Errorf("cannot parse Location header: %w", err)
	}

	err = NewHTTPRequest(r.result, r.client, http.MethodGet, locationURL.Path, nil).Do(ctx)
	if err != nil {
		return fmt.Errorf("cannot follow Location header: %w", err)
	}

	return nil
}

func decompressBody(resp *http.Response) (io.ReadCloser, error) {
	body := resp.Body
	encoding := resp.Header.Get("Content-Encoding")
	switch encoding {
	case "":
		return io.NopCloser(resp.Body), nil
	case "gzip":
		gzipReader, err := gzip.NewReader(body)
		if err != nil {
			return nil, fmt.Errorf("cannot decode gzip response: %w", err)
		}
		return gzipReader, nil
	default:
		return nil, fmt.Errorf(`unexpected Content-Encoding "%s"`, encoding)
	}
}

func decodeBody(resp *http.Response, reader io.ReadCloser, result any) error {
	// Handle special types
	switch result := result.(type) {
	case *NoResult:
		return nil // nop
	case *RawResponse:
		result.Response = resp
		result.Body = reader
		return nil
	}

	// Handle other result types
	contentType, _, _ := strings.Cut(resp.Header.Get("Content-Type"), ";")
	switch contentType {
	case ContentTypeJSON:
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(result); err != nil {
			return fmt.Errorf("cannot decode dataverse response: %w", err)
		}
	default:
		return fmt.Errorf(`unexpected Content-Type "%s"`, contentType)
	}
	return nil
}
