FROM ubuntu:24.04@sha256:2a1e42397521001f21178a06e37ba1024481d3e8b6a754902ac5fb6a0861c7ac

# Install dependencies and create a user with user group
User root

RUN groupadd --gid 1042 vapusdata

RUN useradd -m -u 503 -g vapusdata vapusdata-cli

RUN apt-get update

RUN apt install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

RUN mkdir -p /src

COPY ./vapusctl /src

RUN apt-get update && \
apt-get install -y --no-install-recommends \
     ca-certificates && \
    apt-get clean


RUN apt install --reinstall ca-certificates

RUN update-ca-certificates

WORKDIR /src

RUN chown -R vapusdata-cli:vapusdata .

USER vapusdata-cli

WORKDIR /src

EXPOSE 9798

RUN date -u

CMD ["./vapusctl"]