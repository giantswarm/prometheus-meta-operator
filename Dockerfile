FROM alpine:3.20.0

RUN apk add --no-cache ca-certificates

ADD ./files /opt/prometheus-meta-operator/files

ADD ./prometheus-meta-operator /prometheus-meta-operator

ENTRYPOINT ["/prometheus-meta-operator"]
