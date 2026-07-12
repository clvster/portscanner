FROM golang:1.26.5-alpine3.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/scanner ./cmd/scanner

FROM alpine:3.24
RUN apk add --no-cache masscan libcap && \
    setcap cap_net_raw,cap_net_admin=eip $(which masscan)
COPY --from=build /out/scanner /usr/local/bin/scanner
COPY config /app/config
COPY db/migrations /app/db/migrations
WORKDIR /app
ENTRYPOINT ["scanner"]
