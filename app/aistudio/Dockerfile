FROM ubuntu:24.04@sha256:2a1e42397521001f21178a06e37ba1024481d3e8b6a754902ac5fb6a0861c7ac

# Install dependencies and create a user with user group
USER root

RUN groupadd --gid 1042 vapusdata

RUN useradd -m -u 503 -g vapusdata vapusdata-aistudio-app

RUN apt-get update

RUN apt install -y --no-install-recommends ca-certificates

RUN apt-get install -y --no-install-recommends wget

RUN update-ca-certificates

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.28 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64

RUN chmod +x /bin/grpc_health_probe

RUN apt-get purge -y wget

RUN mkdir -p /src

COPY ./vapus_aistudio /src

WORKDIR /src

RUN chown -R vapusdata-aistudio-app:vapusdata .

USER vapusdata-aistudio-app

EXPOSE 9017

RUN date -u

CMD ["./vapus_aistudio","--conf","/data/vapusdata"]