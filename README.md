# grpcerrors [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](http://godoc.org/github.com/neilgarb/grpcerrors) [![Build Status](https://travis-ci.org/neilgarb/grpcerrors.svg?branch=master)](https://travis-ci.org/neilgarb/grpcerrors)

Register errors so that they can be safely transmitted over gRPC.

In both client and server:

```go
someError := errors.New("some error")
grpcerrors.Register(1, someError)
```

In the server:

```go
return nil, grperrors.ToProto(err)
```

In the client:

```go
_, err := grpcCall()
if grpcerrors.FromProto(err) == someError {
  // ...
}
```
