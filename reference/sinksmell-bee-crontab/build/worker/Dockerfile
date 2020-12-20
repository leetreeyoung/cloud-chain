FROM 0c7swsv7.mirror.aliyuncs.com/library/debian:stretch

COPY bin/worker /root/worker
COPY conf/worker.yaml /worker.yaml

ENTRYPOINT ["/root/worker"]