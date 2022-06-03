FROM golang:alpine AS build

WORKDIR /app/

RUN apk add --no-cache build-base libwebp-dev

COPY . .

RUN  --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags "-s -w" main.go

FROM alpine:edge

RUN apk add --no-cache libwebp

WORKDIR /app/

COPY --from=build /app/main /app/http3-ytproxy

CMD ./http3-ytproxy
