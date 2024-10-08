## Protos

The protobuf definitions are in the `backend/proto` directory. The `backend/Makefile` has a rule for compiling these definitions with `protoc`:

```
cd backend
make protoc
```

The generated files are created in the `backend/protogen/golang` directory.

### gRPC-gateway

The [gRPC-gateway docs](https://grpc-ecosystem.github.io/grpc-gateway/) give a great summary of their usage:
> "gRPC-Gateway is a plugin of protoc. It reads a gRPC service definition and generates a reverse-proxy server which translates a RESTful JSON API into gRPC. This server is generated according to custom options in your gRPC definition."

### Developing protos

The workflow consists of:

1. If adding a new service, add a directory and .proto file with service under the `/backend/proto` directory, RPC, and message definitions; otherwise extend an existing service with new RPCs and any new messages.
2. Run the `protoc` rule of the `Makefile` to generate files.
3. Create a file in `/backend/services` for the service. Define a private struct that implements the methods of the gRPC service interface. Write business logic in this implementation. Define a `New<gRPCServiceName>Svc` func whose function signature returns the interface and whose function body returns the private struct with member fields for any required dependencies.
4. Prepare any dependencies and instantiate (3) in `backend/cmd/api/main.go`, doing dependency injection as needed. Register the service instance on the main gRPC server.
5. Register the gRPC service handler on the gateway's http serve mux.
6. After rebuilding the images for both the `api` and `gateway` containers, the new API route should be callable at the route defined in the grpc-gateway proto annotations.
