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

// ToGrpc converts the given error to a gRPC error by looking the error's code
// up in the registry. If the error hasn't been registered, the error is
// returned unchanged.
func (r *Registry) ToGrpc(err error) error {
	for k, v := range r.errors {
		if v == err {
			return status.Error(k, err.Error())
		}
	}
	return err
}
 
// FromGrpc converts a gRPC error back to the registered error type. If the
// error isn't a gRPC error, or the status code isn't registered, the error
// is returned as is.
func (r *Registry) FromGrpc(err error) error {
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

var defaultRegistry Registry

// Register registers an error code with the default registry.
func Register(code codes.Code, err error) {
	defaultRegistry.Register(code, err)
}

// ToGrpc converts the given error to a gRPC error using the default registry.
func ToGrpc(err error) error {
	return defaultRegistry.ToGrpc(err)
}

// FromGrpc converts a gRPC back to a registered error using the default
// registry.
func FromGrpc(err error) error {
	return defaultRegistry.FromGrpc(err)
}
