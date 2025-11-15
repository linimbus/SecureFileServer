FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o SecureFileServer .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/SecureFileServer .

EXPOSE 8080

CMD ["./SecureFileServer"]