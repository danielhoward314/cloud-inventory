# Backend for cloud-inventory

## API

The API is a gRPC server.

## Gateway

The Google `grpc-gateway` is used to receive RESTful JSON API requests and translate them into gRPC requests that are reverse proxied to the API.

## Protos

The protobuf definitions are in the `proto` directory. The `Makefile` has a rule for compiling these definitions with `protoc`:

```
make protoc
```