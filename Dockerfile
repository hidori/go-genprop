FROM golang:1.25.4-alpine3.22 AS builder

RUN apk update && apk add --no-cache \
    ca-certificates=20250619-r0 \
    make=4.4.1-r3 \
    && update-ca-certificates

COPY . /work
WORKDIR /work

RUN make build

FROM golang:1.25.4-alpine3.22 AS runner

COPY --from=builder /work/bin/genprop /usr/local/bin

ENTRYPOINT [ "genprop" ]
