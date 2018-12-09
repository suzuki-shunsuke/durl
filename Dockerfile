FROM golang:1.11.2 AS go-build-env
ARG DURL_VERSION=0.1.0-1
RUN wget https://github.com/suzuki-shunsuke/durl/releases/download/v${DURL_VERSION}/durl_${DURL_VERSION}_linux_amd64.tar.gz
RUN tar xvzf durl_${DURL_VERSION}_linux_amd64.tar.gz

FROM busybox:1.29.3
COPY --from=go-build-env /go/durl /usr/local/bin/
COPY --from=go-build-env /etc/ssl/certs/ /etc/ssl/certs
