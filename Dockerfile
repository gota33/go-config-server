# This file is a template, and might need editing before it works on your project.
FROM golang:1.16 AS builder

ARG VERSION

WORKDIR /usr/src/server
COPY . .
RUN go build -v -o server -ldflags "-X 'github.com/gota33/go-config-server/internal.Version=$VERSION'" main.go

FROM buildpack-deps:jessie

COPY --from=builder /usr/src/server/server /usr/local/bin
CMD ["server", "web"]
