package protovalidate

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"connectrpc.com/connect"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

type Interceptor struct {
	validator *protovalidate.Validator
	plugins   []InterceptorPlugin
}

func NewValidator() (*protovalidate.Validator, error) {
	return protovalidate.New()
}

func NewInterceptor(validator *protovalidate.Validator, plugins ...InterceptorPlugin) *Interceptor {
	return &Interceptor{validator: validator, plugins: plugins}
}

// WrapUnary implements connect.Interceptor.
func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := i.validate(req.Any()); err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

// WrapStreamingClient implements connect.Interceptor.
func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return &streamingClientInterceptor{StreamingClientConn: next(ctx, spec), interceptor: i}
	}
}

// WrapStreamingHandler implements connect.Interceptor.
func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return next(ctx, &streamingHandlerInterceptor{StreamingHandlerConn: conn, interceptor: i})
	}
}

func (i *Interceptor) validate(msg any) error {
	// Check message type
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("expected proto.Message, got %T", msg)
	}

	// Compose chain of plugins
	var final ValidateFn = i.validator.Validate
	for _, plugin := range slices.Backward(i.plugins) {
		final = plugin(final)
	}

	// Validate with plugins
	err := final(protoMsg)
	if err == nil {
		return nil
	}

	// Add error code
	connectErr := connect.NewError(connect.CodeInvalidArgument, err)

	// Add error details
	var validationErr *protovalidate.ValidationError
	if errors.As(err, &validationErr) {
		if detail, err := connect.NewErrorDetail(validationErr.ToProto()); err == nil {
			connectErr.AddDetail(detail)
		}
	}

	return connectErr
}

type streamingClientInterceptor struct {
	connect.StreamingClientConn
	interceptor *Interceptor
}

func (s *streamingClientInterceptor) Send(msg any) error {
	if err := s.interceptor.validate(msg); err != nil {
		return err
	}
	return s.StreamingClientConn.Send(msg)
}

type streamingHandlerInterceptor struct {
	connect.StreamingHandlerConn
	interceptor *Interceptor
}

func (s *streamingHandlerInterceptor) Receive(msg any) error {
	if err := s.StreamingHandlerConn.Receive(msg); err != nil {
		return err
	}
	return s.interceptor.validate(msg)
}
