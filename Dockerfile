FROM golang as builder

WORKDIR /app

COPY ./broker/go.mod ./broker/go.sum ./

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

RUN export GOPROXY=https://goproxy.cn && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./broker

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/main .

RUN chmod +x /app/main

EXPOSE 8080

ENTRYPOINT ["/app/main"]
