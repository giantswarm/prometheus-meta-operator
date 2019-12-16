FROM alpine:3.10

RUN apk add --no-cache ca-certificates

ADD ./prometheus-meta-operator /prometheus-meta-operator

ENTRYPOINT ["/prometheus-meta-operator"]
