FROM golang:alpine3.12 AS build

WORKDIR /project

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV GO111MODULE=on

RUN go build -o server .


FROM alpine:3.12

WORKDIR /dns-proxy

COPY --from=build /project/server .
COPY config.json ./config/config.json

ENV CONFIG_PATH=/dns-proxy/config/config.json
VOLUME /dns-proxy/config/
EXPOSE 53/udp 80 443

RUN chmod +x /dns-proxy/server

ENTRYPOINT /dns-proxy/server
