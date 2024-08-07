FROM golang:1.22-alpine3.20 AS builder

RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    make \
    && update-ca-certificates

COPY . /workspace
WORKDIR /workspace

RUN make build

FROM alpine:3.20 AS runner

COPY --from=builder /workspace/bin/genprop /usr/local/bin

ENTRYPOINT [ "genprop" ]
