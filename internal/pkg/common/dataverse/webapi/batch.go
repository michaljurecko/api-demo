package webapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type MultipartGenerator func(w *multipart.Writer) error

type MultipartProcessor func(ctx context.Context, r *multipart.Reader) error

type headerProvider interface {
	Get(key string) string
}

func prepareBatchRequest(client *Client, generator MultipartGenerator, onResponse MultipartProcessor) (*HTTPRequest[RawResponse], error) {
	body := &bytes.Buffer{}
	boundary := "batch_" + uuid.New().String()
	writer := multipart.NewWriter(body)
	err := writer.SetBoundary(boundary)
	if err != nil {
		return nil, fmt.Errorf("failed to set multipart boundary: %w", err)
	}

	if err = generator(writer); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req := NewHTTPRequest(&RawResponse{}, client, "POST", "$batch", bytes.NewReader(body.Bytes()))
	req.Header("If-None-Match", "null")
	req.Header("Content-Type", "multipart/mixed; boundary="+boundary)
	req.OnSuccess(func(ctx context.Context, r *RawResponse) error {
		return processMultipartResponse(ctx, r.Response.Header, r.Body, onResponse)
	})

	return req, nil
}

func processMultipartResponse(ctx context.Context, header headerProvider, body io.Reader, onResponse MultipartProcessor) error {
	value := header.Get("Content-Type")
	if len(value) == 0 {
		return errors.New("missing Content-Type header in response")
	}

	// Example: Content-Type: multipart/mixed; boundary=batchresponse_3d1b3e3b-7b6e-4b0b-8f1d-3b6b6b7b6b7b
	contentType := strings.SplitN(value, ";", 2)
	if contentType[0] != "multipart/mixed" {
		return fmt.Errorf("unexpected 'Content-Type' '%s', expected '%s'", contentType, "multipart/mixed")
	}
	if len(contentType) != 2 {
		return errors.New("missing 'boundary' parameter in 'Content-Type' header")
	}

	remainder := strings.TrimSpace(contentType[1])
	if !strings.HasPrefix(remainder, "boundary=") {
		return fmt.Errorf("unexpected parameter in 'Content-Type' header '%s'", remainder)
	}

	batchBoundary := strings.TrimPrefix(remainder, "boundary=")
	batchReader := multipart.NewReader(body, batchBoundary)
	return onResponse(ctx, batchReader)
}

func writePartToBatchRequest(index int, req *http.Request, writer *multipart.Writer) (err error) {
	headers := textproto.MIMEHeader{
		"Content-Type":              {"application/http"},
		"Content-Transfer-Encoding": {"binary"},
		"Content-ID":                {strconv.Itoa(index + 1)},
	}
	if changeSetPart, err := writer.CreatePart(headers); err != nil {
		return fmt.Errorf("failed to create part for request '%d': %w", index+1, err)
	} else if dump, err := httputil.DumpRequest(req, true); err != nil {
		return fmt.Errorf("failed to dump request '%d': %w", index+1, err)
	} else if _, err := changeSetPart.Write(dump); err != nil {
		return fmt.Errorf("failed to write request '%d': %w", index+1, err)
	}
	return nil
}
