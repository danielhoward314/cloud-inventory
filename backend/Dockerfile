FROM golang:1.21 as gobase

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy all .go files
COPY . .

# Build the CLI binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /cloud-inventory-cli ./cmd/cli/main.go

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /cloud-inventory-api ./cmd/api/main.go

# Build the Gateway binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /cloud-inventory-gateway ./cmd/gateway/main.go

############################################
# api
############################################
FROM scratch as api
COPY --from=gobase /cloud-inventory-api /bin/api
EXPOSE 50051
CMD ["/bin/api"]

############################################
# gateway
############################################
FROM scratch as gateway
COPY --from=gobase /cloud-inventory-gateway /bin/gateway
EXPOSE 8080
CMD ["/bin/gateway"]

############################################
# cli
############################################
FROM ubuntu:latest as cli
COPY --from=gobase /cloud-inventory-cli /bin/cli
RUN chmod +x /bin/cli
ENTRYPOINT ["/bin/cli"]