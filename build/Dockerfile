FROM ubuntu:latest
LABEL authors="299m"

## Copy in the config and bin dir
COPY ./config /config
COPY ./bin /bin

ENTRYPOINT ["httpprxy", "-cfgdir", "./config"]