# ============================
# Stage 1: Build app
# ============================

FROM golang:1.27-alpine:3.24.3 as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum ./

COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY order/go.mod order/go.sum ./order/
COPY payment/go.mod payment/go.sum ./payment/

RUN go mod download

COPY platform ./platform
COPY shared ./shared
COPY inventory ./inventory
COPY order ./order
COPY payment ./payment

ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

RUN chmod +x grpc-health-probe

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd/grpc_server/main.go

# ============================
# Stage 2: Final image
# ============================

FROM alpine:3.24.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-inventory .

COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

COPY deploy/compose/payment/.env .env

EXPOSE 50052

ENTRYPOINT ["./app-inventory"]
