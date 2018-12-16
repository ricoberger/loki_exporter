FROM   alpine:latest
LABEL  maintainer="Rico Berger <mail@ricoberger.de>"

COPY ./bin/loki_exporter-*-linux-amd64/loki_exporter  /bin/loki_exporter
COPY config.yml     /etc/loki_exporter/config.yml

EXPOSE      9524
ENTRYPOINT  [ "/bin/loki_exporter" ]
CMD         [ "-config.file=/etc/loki_exporter/config.yml" ]
