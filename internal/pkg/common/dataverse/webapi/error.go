package webapi

import "fmt"

type APIError struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UnexpectedStatusError struct {
	Method   string
	URL      string
	Expected int
	Actual   int
	Message  string
}

func (e UnexpectedStatusError) Error() string {
	out := fmt.Sprintf("HTTP request %s '%s' failed: unexpected status %d, expected %d", e.Method, e.URL, e.Actual, e.Expected)
	if e.Message != "" {
		out += ": " + e.Message
	}
	return out
}
