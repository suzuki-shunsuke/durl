FROM linuxkit/ca-certificates:v0.8-amd64 AS ssl

FROM busybox:1.32.0
COPY dist/durl_linux_amd64/durl /usr/local/bin/
COPY --from=ssl /etc/ssl/certs/ /etc/ssl/certs/
