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

In order to create a Cloud Inventory account, navigate to `cloud-inventory-ui.local/signup`.

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

## maildev SMTP mock server

The signup flow for this application requires new users verify their email with a code sent via email. In local development, I use [maildev](https://github.com/maildev/maildev) as a mock SMTP server. In addition to the SMPT server, the `maildev` container also spins up a UI at `http://0.0.0.0:1080/`.

## redis

The `redis` container is used during the signup process for email verification data and is used for user session JWTs. For troubleshooting the `redis` container, you can exec into it and run the Redis CLI with these commands:

```
docker exec -it redis sh
redis-cli
```

Some common commands are:

```
KEYS * # get all keys by a pattern (wildcard used here)
GET <key> # read the data at a given key
FLUSHDB # delete all data
```

## tailwindcss

The Vue SPA uses [tailwindcss](https://tailwindcss.com/docs/installation) for its styling. To install and configure it initially, I did the following:

```
cd frontend/cloud-inventory-spa
npm install -D tailwindcss
npx tailwindcss init
touch ./src/tailwind-input.css ./src/assets/tailwind-output.css
```

I modified the `content` and `theme` properties within the `tailwind.config.js` file to point to the files that use CSS and to use a font family. I also added the following to the `tailwind-input.css` file:

```
@tailwind base;
@tailwind components;
@tailwind utilities;

@import url("https://fonts.googleapis.com/css2?family=Poppins:wght@400;500&display=swap");

@layer base {
    html {
      font-family: "Poppins", system-ui, sans-serif;
    }
}
```

tailwind uses utility classes and only ever generates the CSS for them when they are used. In order to generate them on the fly during development as you use more or fewer utility classes, run this command:

```
npx tailwindcss -i ./src/tailwind-input.css -o ./src/assets/tailwind-output.css --watch
```
