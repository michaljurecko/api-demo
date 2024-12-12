package webapi

import (
	"context"
	"slices"
)

type APIRequest[T any] struct {
	result   *T
	requests []httpRequest
}

type apiRequest interface {
	DoOrErr(ctx context.Context) error
	httpRequests() []httpRequest
}

func NewAPIRequest[T any](result *T, requests ...httpRequest) *APIRequest[T] {
	return &APIRequest[T]{
		result:   result,
		requests: requests,
	}
}

func NewAPIRequestError[T any](result *T, err error) *APIRequest[T] {
	return &APIRequest[T]{
		result:   result,
		requests: []httpRequest{NewHTTPRequestError(err)},
	}
}

func (r *APIRequest[T]) Do(ctx context.Context) (*T, error) {
	for _, req := range r.requests {
		if err := req.Do(ctx); err != nil {
			return nil, err
		}
	}
	return r.result, nil
}

func (r *APIRequest[T]) DoOrErr(ctx context.Context) error {
	_, err := r.Do(ctx)
	return err
}

func (r *APIRequest[T]) httpRequests() []httpRequest {
	return slices.Clone(r.requests)
}
