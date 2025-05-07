FROM ubuntu:24.04

# Install dependencies and create a user with user group
USER root

RUN groupadd --gid 1042 vapusdata

RUN useradd -m -u 503 -g vapusdata vapusdata-vapusai-app

RUN mkdir -p /src

COPY ./dist/aitools/aitools /src

WORKDIR /src

RUN chown -R vapusdata-vapusai-app:vapusdata .

# RUN GRPC_HEALTH_PROBE_VERSION=v0.4.28 && \
#     curl -Lo /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
#     chmod +x /bin/grpc_health_probe

USER vapusdata-vapusai-app

WORKDIR /src

EXPOSE 9025

RUN date -u

CMD ["./aitools","--conf","/data/vapusdata"]