package grpcerrors

import "errors"

type grpcError struct {
	code int
	err error
}

func (e *grpcError) Error() string {
	return e.err.Error()
}

type Registry struct {
	errors map[int]error
}

func NewRegistry() *Registry {
	return &Registry{errors: make(map[int]error)}
}

func (r *Registry) Register(code int, err error) error {
	_, ok := r.errors[code]
	if ok {
		return errors.New("grpcerrors: code already registered")
	}
	r.errors[code] = err
	return nil
}

func (r *Registry) ToGrpc(err error) error {
	for k, v := range r.errors {
		if v == err {
			return &grpcError{code: k, err: err}
		}
	}
	return err
}

func (r *Registry) FromGrpc(err error) error {
	if err == nil {
		return nil
	}
	gerr, ok := err.(*grpcError)
	if !ok {
		return err
	}
	rerr, ok := r.errors[gerr.code]
	if !ok {
		return err
	}
	return rerr
}

var defaultRegistry Registry

func Register(code int, err error) {
	defaultRegistry.Register(code, err)
}

func ToGrpc(err error) error {
	return defaultRegistry.ToGrpc(err)
}

func FromGrpc(err error) error {
	return defaultRegistry.FromGrpc(err)
}
