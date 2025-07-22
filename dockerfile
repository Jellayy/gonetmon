FROM golang:1.24.2-bookworm AS builder

WORKDIR /gonetmon

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o server

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /gonetmon/server /gonetmon/server
COPY --from=builder /gonetmon/config.yaml /config.yaml

CMD ["/gonetmon/server"]
