package validation

import (
	"buf.build/go/protovalidate"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BuildValidationError(err error) error {
	vErr, ok := err.(*protovalidate.ValidationError)
	if !ok {
		return status.Error(codes.InvalidArgument, "invalid input")
	}

	st := status.New(codes.InvalidArgument, "validation error")

	br := &errdetails.BadRequest{}

	for _, v := range vErr.Violations {

		field := v.Proto.GetField()
		if field == nil || len(field.Elements) == 0 {
			continue
		}

		name := field.Elements[0].GetFieldName()

		msg := v.Proto.GetMessage()
		if msg == "" {
			msg = "Invalid value"
		}

		// custom message
		if name == "birthday" {
			msg = "birthday must be in format YYYY-MM-DD and a valid date"
		}

		br.FieldViolations = append(
			br.FieldViolations,
			&errdetails.BadRequest_FieldViolation{
				Field:       name,
				Description: msg,
			},
		)
	}

	st, err = st.WithDetails(br)
	if err != nil {
		return status.Error(codes.Internal, "failed to build validation error details")
	}

	return st.Err()
}

func BuildServiceUnavailableError(service string) error {
	st := status.New(
		codes.Unavailable,
		"dependent "+service+" service unavailable",
	)

	retry := &errdetails.RetryInfo{}

	st, err := st.WithDetails(retry)
	if err != nil {
		return status.Error(
			codes.Unavailable,
			service+" service unavailable",
		)
	}

	return st.Err()
}

func BuildBusinessError(code, msg string) error {
	st := status.New(codes.FailedPrecondition, msg)

	info := &errdetails.ErrorInfo{
		Reason: code,
	}

	st, _ = st.WithDetails(info)

	return st.Err()
}
