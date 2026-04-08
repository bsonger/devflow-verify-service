FROM registry.cn-hangzhou.aliyuncs.com/devflow/golang:1.25.8 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOROOT=$(go env GOROOT) swag init -g cmd/main.go --parseDependency -o docs/generated/swagger
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o devflow-verify-service ./cmd

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/devflow-verify-service ./devflow-verify-service
COPY --from=builder /app/docs ./docs

RUN adduser -D devuser
USER devuser

EXPOSE 8080

ENTRYPOINT ["./devflow-verify-service"]
