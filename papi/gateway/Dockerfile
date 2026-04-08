# Stage 1 — compile the plugin shared library.
#
# The KrakenD 2.5 binary was compiled with Go 1.20.13. Go's plugin loader
# verifies that the plugin's runtime package hash matches the host binary's,
# so we MUST compile the plugin with the exact same Go version.
#
# Alpine 3.18 only ships Go 1.20.11 (wrong version). The official Go 1.20.13
# tarball from go.dev is glibc-linked, but KrakenD runs on Alpine (musl).
# gcompat provides a glibc→musl shim so the Go 1.20.13 toolchain can run here.
# CGO then uses Alpine's musl-targeting gcc to produce a musl-compatible .so.
FROM devopsfaith/krakend:2.5 AS plugin-builder

ARG TARGETARCH=arm64
ARG GO_VERSION=1.20.13

RUN apk add --no-cache musl-dev gcc binutils binutils-gold wget ca-certificates gcompat

RUN wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" \
         -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go

WORKDIR /plugin
COPY plugins/error-status-propagator/ .
RUN go version && \
    CGO_ENABLED=1 CC="gcc -fuse-ld=gold" go build -buildmode=plugin -o /tmp/error-status-propagator.so .

# Stage 2 — KrakenD with the compiled plugin
FROM devopsfaith/krakend:2.5

COPY --from=plugin-builder /tmp/error-status-propagator.so /opt/krakend/plugins/
