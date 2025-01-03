package protovalidate

import "google.golang.org/protobuf/proto"

// InterceptorPlugin provides modification of validation results.
type InterceptorPlugin func(next ValidateFn) ValidateFn

type ValidateFn func(msg proto.Message) error
