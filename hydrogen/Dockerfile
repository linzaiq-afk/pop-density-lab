ARG H2O_REPO=https://github.com/h2o/h2o.git
ARG H2O_REF=master

# --- build stage -----------------------------------------------------------
FROM ubuntu:24.04 AS build

ARG H2O_REPO
ARG H2O_REF

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        cmake \
        make \
        g++ \
        pkg-config \
        libssl-dev \
        zlib1g-dev \
        libbrotli-dev \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --depth 1 --branch "${H2O_REF}" "${H2O_REPO}" /src/h2o \
    && cd /src/h2o \
    && git submodule update --init --recursive

RUN cd /src/h2o \
    && cmake -DCMAKE_BUILD_TYPE=Release -DWITH_MRUBY=OFF . \
    && make -j"$(nproc)" \
    && make install

# --- runtime stage ---------------------------------------------------------
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        libssl3t64 \
        libbrotli1 \
        zlib1g \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /usr/local/bin/h2o /usr/local/bin/h2o
COPY h2o.conf /etc/hydrogen/h2o.conf

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/h2o"]
CMD ["-c", "/etc/hydrogen/h2o.conf"]
