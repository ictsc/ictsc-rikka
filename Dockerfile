### ビルド用コンテナ ###
FROM golang:1.19 AS builder

WORKDIR /go/src/ictsc-rikka/

COPY . .

# go build -o rikka cmd/rikka/*.go
RUN make build


### メインコンテナ ###
FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /go/src/ictsc-rikka/rikka .

EXPOSE 8080

CMD [ "./rikka" ]
