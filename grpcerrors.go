// Package grpcerrors helps transmit errors over gRPC by associating errors
// with status codes which can be transmitted safely.
package grpcerrors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Registry holds mappings from error codes to errors.
type Registry struct {
	errors map[codes.Code]error
}

// NewRegistry returns a new error code registry.
func NewRegistry() *Registry {
	return &Registry{errors: make(map[codes.Code]error)}
}

// Register associates a status code with an error.
func (r *Registry) Register(code codes.Code, err error) error {
	_, ok := r.errors[code]
	if ok {
		return errors.New("grpcerrors: code already registered")
	}
	r.errors[code] = err
	return nil
}

// MustRegister registers an error code with the registry and panics
// if the code as already registered.
func (r *Registry) MustRegister(code codes.Code, err error) {
	if err := r.Register(code, err); err != nil {
		panic(err)
	}
}

// ToProto converts the given error to a gRPC error by looking the error's code
// up in the registry. If the error hasn't been registered, the error is
// returned unchanged.
func (r *Registry) ToProto(err error) error {
	for k, v := range r.errors {
		if v == err {
			return status.Error(k, err.Error())
		}
	}
	return err
}

// FromProto converts a gRPC error back to the registered error type. If the
// error isn't a gRPC error, or the status code isn't registered, the error
// is returned as is.
func (r *Registry) FromProto(err error) error {
	if err == nil {
		return nil
	}
	s, ok := status.FromError(err)
	if !ok {
		return err
	}
	rerr, ok := r.errors[s.Code()]
	if !ok {
		return err
	}
	return rerr
}

var defaultRegistry = NewRegistry()

// Register registers an error code with the default registry.
func Register(code codes.Code, err error) error {
	return defaultRegistry.Register(code, err)
}

// MustRegister registers an error code with the default registry and panics
// if the code as already registered.
func MustRegister(code codes.Code, err error) {
	defaultRegistry.MustRegister(code, err)
}

// ToProto converts the given error to a gRPC error using the default registry.
func ToProto(err error) error {
	return defaultRegistry.ToProto(err)
}

// FromProto converts a gRPC back to a registered error using the default
// registry.
func FromProto(err error) error {
	return defaultRegistry.FromProto(err)
}
