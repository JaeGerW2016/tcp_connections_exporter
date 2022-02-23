FROM alpine:3.14.1

EXPOSE 9139

COPY tcp_connections_exporter /usr/bin/tcp_connections_exporter
ENTRYPOINT [ "/usr/bin/tcp_connections_exporter" ]