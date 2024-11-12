FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main main.go 

FROM alpine:3.20.3

WORKDIR /app

COPY --from=builder /app/main .

USER 1001

CMD ["./main"]