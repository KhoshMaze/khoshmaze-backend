FROM golang:1.23 as builder 

WORKDIR /github.com/KhoshMaze/khoshmaze-backend 

RUN CGO_ENABLED=0 go build -o ./khoshmaze-api cmd/main.go

FROM alpine as deploy

RUN apk add --no-cache tzdata
ENV TZ=Asia/Tehran

WORKDIR /github.com/KhoshMaze/khoshmaze-backend 

COPY --from=builder /github.com/KhoshMaze/khoshmaze-backend/khoshmaze-api ./khoshmaze-api

CMD ["./khoshmaze-api" , "--config" , "/etc/config.json"]