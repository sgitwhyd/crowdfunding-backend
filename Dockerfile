FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.* /app/

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 3001

# Run the app
CMD ["./main"]
