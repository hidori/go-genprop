FROM golang:1.25-alpine3.22 AS builder

RUN apk update && apk add --no-cache \
    ca-certificates \
    make \
    && update-ca-certificates

COPY . /work
WORKDIR /work

RUN make build

FROM golang:1.25-alpine3.22 AS runner

COPY --from=builder /work/bin/genprop /usr/local/bin

ENTRYPOINT [ "genprop" ]
