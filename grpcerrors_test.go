package grpcerrors_test

import (
	"errors"
	"testing"

	"github.com/neilgarb/grpcerrors"
)

func TestRegistry(t *testing.T) {
	r := grpcerrors.NewRegistry()

	e1 := errors.New("foo")
	if err := r.Register(1, e1); err != nil {
		t.Error("Expected success, got:", err)
	}
	e2 := errors.New("bar")
	if err := r.Register(1, e2); err == nil {
		t.Error("Expected error, got success")
	}
	if err := r.Register(2, e2); err != nil {
		t.Error("Expected success, got:", err)
	}

	e11 := r.ToGrpc(e1)
	if e11 == e1 {
		t.Errorf("Expected different errors")
	}
	e12 := r.FromGrpc(e11)
	if e12 != e1 {
		t.Error("Expected same error, got:", e12, e1)
	}
	e13 := r.FromGrpc(e1)
	if e13 != e1 {
		t.Error("Expected same error, got:", e13, e1)
	}
}
