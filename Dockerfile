FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /mysql-waiter .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /mysql-waiter /mysql-waiter

EXPOSE 3000
HEALTHCHECK --interval=1s --timeout=1s --retries=1 --start-interval=1s --start-period=24h CMD ["/mysql-waiter", "-health"]
ENTRYPOINT ["/mysql-waiter"]
