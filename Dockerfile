FROM golang:1.23-alpine3.20 AS builder

RUN apk update && apk add --no-cache \
    ca-certificates \
    make \
    && update-ca-certificates

COPY . /work
WORKDIR /work

RUN make build

FROM alpine:3.20 AS runner

COPY --from=builder /work/bin/genprop /usr/local/bin

ENTRYPOINT [ "genprop" ]
