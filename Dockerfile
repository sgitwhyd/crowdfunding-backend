FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:3.18

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy the binary and configs from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

RUN chmod +x /app/main

EXPOSE 3001

# Run the app
CMD ["./main"]
