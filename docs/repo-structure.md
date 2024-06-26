# Repo Structure

The `backend` directory contains three Go binaries:

1. The binary in `backend/cmd/api` is the gRPC server code that handles the business logic of the API.
2. The binary in `backend/cmd/gateway` uses the Google grpc-gateway to translate JSON RESTful API requests to protobufs and reverse proxies them to the API in (1).
3. The binary in `backend/cmd/cli` is a CLI tool for managing the application database and SQL migrations for its tables.

The `frontend` directory contains the Vue SPA for the Cloud Inventory UI. This SPA was bootstrapped with the following:

```
npm create vue@latest
```