FROM linuxkit/ca-certificates:v0.7-amd64 AS ssl

FROM busybox:1.31.1
COPY dist/durl_linux_amd64/durl /usr/local/bin/
COPY --from=ssl /etc/ssl/certs/ /etc/ssl/certs/
