
FROM golang:1.22.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod vendor
COPY . .
RUN chmod +x scripts/generate.sh
RUN ./scripts/generate.sh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/my-go-app .
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/my-go-app /app/my-go-app
# Expose the port your Go app listens on
EXPOSE 8080
# Command to run the Go application
CMD ["/app/my-go-app"]