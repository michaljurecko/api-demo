package webapi

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sync"

	"github.com/google/uuid"
)

type ChangeSet struct {
	client   *Client
	requests []httpRequest
}

// changeSetRequest is an auxiliary struct with prepared raw requests.
type changeSetRequest struct {
	client      *Client
	requests    []httpRequest
	rawRequests []*http.Request
}

func NewChangeSet(client *Client) *ChangeSet {
	return &ChangeSet{client: client}
}

func (s *ChangeSet) Add(request apiRequest) (contentID int) {
	httpRequests := request.httpRequests()
	if len(httpRequests) != 1 {
		panic(errors.New("cannot add request to ChangeSet: only API requests wrapping one HTTP request are supported"))
	}

	s.requests = append(s.requests, httpRequests[0])
	return len(s.requests) // 1,2,3 ...
}

func (s *ChangeSet) Do(ctx context.Context) error {
	req, err := s.prepareRequest(ctx)
	if err != nil {
		return err
	}
	return req.Do(ctx)
}

func (s *ChangeSet) prepareRequest(ctx context.Context) (*changeSetRequest, error) {
	req := &changeSetRequest{client: s.client, requests: s.requests}

	// Prepare raw HTTP requests
	for _, spec := range s.requests {
		httpReq, err := spec.prepareRequest(ctx)
		if err != nil {
			return nil, err
		}

		// Content-ID references don't work in batch requests,
		// if the return=representation preference is used.
		httpReq.Header.Set("Prefer", "return=minimal")

		req.rawRequests = append(req.rawRequests, httpReq)
	}

	return req, nil
}

func (r *changeSetRequest) Do(ctx context.Context) error {
	req, err := prepareBatchRequest(r.client, r.prepareBatchRequest, r.processBatchResponse)
	if err != nil {
		return err
	}
	return req.Do(ctx)
}

func (r *changeSetRequest) prepareBatchRequest(batchWriter *multipart.Writer) (err error) {
	// Create a batch part for the change set
	changeSetBoundary := "changeset_" + uuid.New().String()
	batchPart, err := batchWriter.CreatePart(textproto.MIMEHeader{"Content-Type": {"multipart/mixed; boundary=" + changeSetBoundary}})
	if err != nil {
		return fmt.Errorf("failed to create batch part: %w", err)
	}

	// Prepare the change set writer
	changeSetWriter := multipart.NewWriter(batchPart)
	err = changeSetWriter.SetBoundary(changeSetBoundary)
	if err != nil {
		return fmt.Errorf("failed to set multipart boundary: %w", err)
	}

	// Compose individual requests to the batch request
	for i, req := range r.rawRequests {
		if err := writePartToBatchRequest(i, req, changeSetWriter); err != nil {
			return err
		}
	}

	if err = changeSetWriter.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return nil
}

func (r *changeSetRequest) processBatchResponse(ctx context.Context, batchReader *multipart.Reader) (err error) {
	// Parse the batch response
	changeSetPart, err := batchReader.NextPart()
	if err != nil {
		return fmt.Errorf("failed to get next part: %w", err)
	}
	defer func() {
		if closeErr := changeSetPart.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	// Iterate over each change set response part
	return processMultipartResponse(ctx, changeSetPart.Header, changeSetPart, r.processBatchResponseParts)
}

func (r *changeSetRequest) processBatchResponseParts(ctx context.Context, batchReader *multipart.Reader) (err error) {
	var index int

	var wg sync.WaitGroup
	var lock sync.Mutex
	var errs []error

	for {
		// Each response part is a response for an original request part
		part, err := batchReader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to get next part: %w", err)
		}

		// Read part body
		partBody, err := io.ReadAll(part)
		if err != nil {
			return fmt.Errorf("failed to read response part '%d': %w", index+1, err)
		}

		// Close part body
		err = part.Close()
		if err != nil {
			return fmt.Errorf("failed to close response part '%d': %w", index+1, err)
		}

		// Process part body as usual
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			if err := r.processBatchResponseOnePart(ctx, index, partBody); err != nil {
				lock.Lock()
				errs = append(errs, err)
				lock.Unlock()
			}
		}(index)

		index++
	}

	wg.Wait()
	return errors.Join(errs...)
}

func (r *changeSetRequest) processBatchResponseOnePart(ctx context.Context, index int, partBody []byte) error {
	// Convert part to the *http.Response
	req := r.requests[index]
	rawReq := r.rawRequests[index]
	if resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(partBody)), rawReq); err != nil {
		return fmt.Errorf("failed to read response part '%d': %w", index+1, err)
	} else if err = req.processResponse(ctx, resp); err != nil {
		return fmt.Errorf("failed to process response part '%d': %w", index+1, err)
	} else if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("failed to close response part '%d': %w", index+1, err)
	}
	return nil
}
