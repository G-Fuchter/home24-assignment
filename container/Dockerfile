FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o build/server ./cmd/server.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/server ./

COPY --from=builder /app/static ./static

EXPOSE 8080

ENTRYPOINT ["./server"]
