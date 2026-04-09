
FROM golang:1.26 AS builder

COPY . /src
WORKDIR /src


RUN go env -w GOPROXY=direct
RUN go mod download


RUN go build -o /app/traceability ./cmd/traceability

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
	gettext \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

WORKDIR /app


COPY --from=builder /app/traceability /app/traceability


COPY configs /app/configs

EXPOSE 8000
EXPOSE 9000


CMD ["sh", "-c", "envsubst < /app/configs/config_docker.yaml > /app/configs/config_resolved.yaml && /app/traceability -conf /app/configs/config_resolved.yaml"]
