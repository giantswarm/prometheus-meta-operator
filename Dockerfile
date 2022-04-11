FROM alpine:3.15.4

RUN apk add --no-cache ca-certificates

ADD ./files /opt/prometheus-meta-operator/files

ADD ./prometheus-meta-operator /prometheus-meta-operator

ENTRYPOINT ["/prometheus-meta-operator"]
