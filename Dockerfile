FROM golang:1.23 AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./khoshmaze-api ./cmd/main.go

FROM alpine AS deploy

WORKDIR /app 

RUN apk add --no-cache tzdata
ENV TZ=Asia/Tehran

COPY --from=builder /app/khoshmaze-api .

# CMD ["./khoshmaze-api", "--config", "/etc/config.json"]