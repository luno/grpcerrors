# grpcerrors

Register errors so that they can be safely transmitted over gRPC.

In both client and server:

```go
someError := errors.New("some error")
grpcerrors.Register(1, someError)
```

In the server:

```go
return nil, grperrors.ToGrpc(err)
```

In the client:

```go
_, err := grpcCall()
if grpcerrors.FromGrpc(err) == someError {
  // ...
}
```
