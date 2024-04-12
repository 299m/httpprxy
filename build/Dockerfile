FROM ubuntu:23.10
LABEL authors="299m"

## Copy in the config and bin dir
COPY ./config /config
COPY ./bin /bin

ENTRYPOINT ["./bin/httpprxy", "-cfgdir", "./config"]