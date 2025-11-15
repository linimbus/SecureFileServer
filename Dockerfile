FROM ubuntu:25.04 AS builder

RUN apt-get update && apt-get install -y \
    build-essential \
    libsqlite3-dev \
    pkg-config \
    golang \
    ca-certificates

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux \
    go build -a -installsuffix cgo -o SecureFileServer .

FROM ubuntu:25.04

WORKDIR /app/

COPY --from=builder /app/SecureFileServer .

RUN chmod +x /app/SecureFileServer

EXPOSE 8080

ENTRYPOINT ["/app/SecureFileServer"]