FROM golang:1.20.2-alpine3.17 as deps

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM deps as builder

WORKDIR /app

RUN go build -o main .

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/main .

# Command expects that the config file is mounted to /app/server.yaml
CMD ["/app/main"]