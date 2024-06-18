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
cd frontend/cloud-inventory-spa && npm run dev
```

The API listens at TCP address `[::]:50051`, the gateway listens at `http://localhost:8080`, and the Vue dev server serves the SPA at `http://localhost:5173`.

## Running on Docker containers

A prerequisite is to modify your hosts file so some hostnames will resolve to the loopback address:

```
sudo vi /etc/hosts

# append these lines
127.0.0.1 cloud-inventory-ui.local
::1 cloud-inventory-ui.local
127.0.0.1 cloud-inventory-gateway.local
::1 cloud-inventory-gateway.local
```

Use docker compose to build the `api`, `gateway`, and `ui` containers:

```
docker compose --file compose.yml up
```