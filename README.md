# Glofox API Client

This is a Golang client for the Glofox API, providing endpoints for managing bookings and classes:

- **Classes**
  - `GET`, `POST`
- **Bookings**
  - `GET`, `POST`

## Getting Started

You can build and run the service using Docker.

### Build and Run

```bash
docker-compose up --build
```

To View Swagger Docs (OPENAPI SPEC) after docker is up and running visit
http://localhost:8080/swagger/index.html

### Test

```bash
go test -v ./...
```
To run only golang server locally without mongoDB and docker
you can run following commands

```bash
go mod tidy
```

```bash
go mod vendor
```

```bash
go run main.go
```
PreRequisite: you need to have new mongodb server up and running at port 27017
