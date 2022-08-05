FROM golang:1.18.4-alpine as build-stage

RUN apk --no-cache add \
    g++ \
    git \
    make \
    bash

WORKDIR ../../
COPY . .
RUN go mod download
RUN go build main.go && \
    chmod 770 main && \
    mv main loadbalancerip-mutator

# Final image.
FROM alpine:latest
RUN apk --no-cache add \
    ca-certificates
COPY --from=build-stage loadbalancerip-mutator /usr/local/bin/loadbalancerip-mutator
ENTRYPOINT ["/usr/local/bin/loadbalancerip-mutator"]