FROM golang:alpine AS build

WORKDIR /app/

COPY . .

RUN go build -ldflags "-s -w" main.go

FROM alpine:edge

WORKDIR /app/

COPY --from=build /app/main /app/http3-ytproxy

CMD ./http3-ytproxy
