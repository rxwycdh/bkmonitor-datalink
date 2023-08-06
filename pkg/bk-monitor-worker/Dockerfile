FROM centos:7

LABEL maintainer="Tencent BlueKing"

RUN mkdir -p /data/bkmonitor/logs  /data/bkmonitor/config

WORKDIR /data/bkmonitor

ADD bin/bmw /data/bkmonitor/bmw
RUN chmod +x /data/bkmonitor/bmw

ENTRYPOINT ["/data/bkmonitor/bmw"]
CMD ["-c", "/data/bkmonitor/config/bmw.yaml"]