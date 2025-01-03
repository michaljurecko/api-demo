package protovalidate

import (
	"errors"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const UpdateMaskFieldName = protoreflect.Name("update_mask")

func UpdateMaskPlugin(next ValidateFn) ValidateFn {
	return func(msg proto.Message) error {
		err := next(msg)

		maskField := updateMaskField(msg)
		if maskField == nil {
			return err
		}

		var validationErr *protovalidate.ValidationError
		if !errors.As(err, &validationErr) {
			return err
		}

		activeFields := make(map[protoreflect.Name]struct{})
		for _, field := range maskField.GetPaths() {
			activeFields[protoreflect.Name(field)] = struct{}{}
		}

		var filteredViolations []*protovalidate.Violation
		for _, v := range validationErr.Violations {
			if _, ok := activeFields[v.FieldDescriptor.Name()]; ok {
				filteredViolations = append(filteredViolations, v)
			}
		}

		validationErr.Violations = filteredViolations

		if len(validationErr.Violations) == 0 {
			return nil
		}

		return err
	}
}

func updateMaskField(msg any) *fieldmaskpb.FieldMask {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil
	}

	// Is method name Update*?
	typeName := string(protoMsg.ProtoReflect().Descriptor().Name())
	if !strings.HasPrefix(typeName, "Update") {
		return nil
	}

	// Is there an FieldMask field?
	var mask *fieldmaskpb.FieldMask
	protoMsg.ProtoReflect().Range(func(field protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		if field.Name() == UpdateMaskFieldName {
			if v, ok := value.Message().Interface().(*fieldmaskpb.FieldMask); ok {
				mask = v
				return false
			}
		}
		return true
	})
	return mask
}
