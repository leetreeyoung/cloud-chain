FROM 0c7swsv7.mirror.aliyuncs.com/library/debian:stretch

COPY bin/master /root/master
COPY conf /root/conf
COPY static /root/static
COPY views /root/views
COPY swagger /root/swagger

CMD ["bash","-c","/root/master","&",">/root/master.log"]