# Cloud Inventory

## Design

TODO: link a png with diagram.
TODO: link the blog.

## Repo Structure

The `backend` directory contains two Go binaries:

1. The binary in `backend/cmd/api` is the gRPC server code that handles the business logic of the API.
2. The binary in `backend/cmd/api` uses the Google grpc-gateway to translate JSON RESTful API requests to protobufs and reverse proxies them to the API in (1).

The `frontend` directory contains the Vue SPA for the Cloud Inventory UI.

## Running locally

Requires 3 separate terminal tabs:

```
go run backend/cmd/api/main.go
go run backend/cmd/gateway/main.go
npm run dev
```

The API listens at TCP address `[::]:50051`, the gateway listens at `http://localhost:8080`, and the Vue dev server serves the SPA at `http://localhost:5173`.