# Cloud Inventory

I created Cloud Inventory to have a fully-featured SaaS product in my portfolio of personal projects.

Cloud Inventory provides dashboards for viewing your cloud resources on providers like AWS, GCP, and Azure.

For more details, read my dev log [here](https://danielhoward-dev.netlify.app/).

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

In order to create a Cloud Inventory account, navigate to `cloud-inventory-ui.local/signup`. For return visits, log in at `cloud-inventory-ui.local/login`.

## Documentation

- [Architecture](./docs/architecture.md)
- [Repo Structure](./docs/repo-structure.md)
- [DB Migrations](./docs/db-migrations.md)
- [Protos](./docs/protos.md)
- [Local containers](./docs/local-containers.md) such as `redis` or `maildev`
- [tailwindcss](./docs/tailwindcss.md)
- [authorization](./docs/authorization.md)