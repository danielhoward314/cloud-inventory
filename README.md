# Cloud Inventory

## Design

TODO: link a png with diagram.
TODO: link the blog.

## Repo Structure

The `backend` directory contains three Go binaries:

1. The binary in `backend/cmd/api` is the gRPC server code that handles the business logic of the API.
2. The binary in `backend/cmd/gateway` uses the Google grpc-gateway to translate JSON RESTful API requests to protobufs and reverse proxies them to the API in (1).
3. The binary in `backend/cmd/cli` is a CLI tool for managing the application database and SQL migrations for its tables.

The `frontend` directory contains the Vue SPA for the Cloud Inventory UI. This SPA was bootstrapped with the following:

```
npm create vue@latest
```

## Running with docker compose

A prerequisite is to modify your hosts file so some hostnames will resolve to the loopback address:

```
sudo vi /etc/hosts

# append these lines
127.0.0.1 cloud-inventory-ui.local
::1 cloud-inventory-ui.local
127.0.0.1 cloud-inventory-gateway.local
::1 cloud-inventory-gateway.local
```

Use docker compose to build the `cli`, `api`, `gateway`, and `ui` containers:

```
docker compose --file compose.yml build
```

The `ui` container is an nginx server that serves the vite production build of the Vue SPA as static files.

## DB Migrations

The `cli` container is a [cobra CLI](https://github.com/spf13/cobra) that exposes commands for initializing the application database. The CLI exposes additional commands that leverage [goose](https://github.com/pressly/goose) for executing the SQL migration files against the application database.

Build the `cli`:

```
docker compose build cli
```

Run the `postgres` container:

```
docker compose up -d postgres
```

To create the main application database:

```
docker compose run --rm cli create db
```

This will create a database with the name matching the `POSTGRES_APPLICATION_DATABASE` value in `env/postgres`.

To drop the main application database:

```
docker compose run --rm cli drop db
```

To create a SQL migration file:

```
docker compose run --rm cli create migration <name>
```

This will create an empty migration file of the format `<timestamp>_<name>.sql` where `<name>` corresponds to the given CLI argument. The file will be in the `backend/cmd/cli/commands/migrations` directory.

To run all migrations through goose:

```
docker compose run --rm cli migrate up
```

The `goose` tool maintains its own table tracking which migrations have been run. The `up` command will only run migrations that have yet to be run, allowing the lifecycle of the database schemas to exist in version control.

To roll back all migrations through goose:

```
docker compose run --rm cli migrate down
```

For troubleshooting the `postgres` container, you can exec into it with this command:

```
docker exec -it postgres psql -U postgres postgres
```

## Protos

The protobuf definitions are in the `backend/proto` directory. The `backend/Makefile` has a rule for compiling these definitions with `protoc`:

```
make protoc
```

The generated files are created in the `backend/protogen/golang` directory.

### gRPC-gateway

The [gRPC-gateway docs](https://grpc-ecosystem.github.io/grpc-gateway/) give a great summary of their usage:
> "gRPC-Gateway is a plugin of protoc. It reads a gRPC service definition and generates a reverse-proxy server which translates a RESTful JSON API into gRPC. This server is generated according to custom options in your gRPC definition."
