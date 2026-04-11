FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
        gettext \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

WORKDIR /app

# Binary and configs are pre-built by 'make build' before docker build
COPY bin/traceability /app/traceability
COPY configs /app/configs

EXPOSE 8000
EXPOSE 9000

CMD ["sh", "-c", "envsubst < /app/configs/config_docker.yaml > /tmp/config_resolved.yaml && /app/traceability -conf /tmp/config_resolved.yaml"]
