FROM linuxkit/ca-certificates:v0.6-amd64 AS ssl

FROM busybox:1.29.3
COPY dist/linux_amd64/durl /usr/local/bin/
COPY --from=ssl /etc/ssl/certs/ /etc/ssl/certs/
