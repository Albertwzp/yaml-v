FROM scratch

ENV TZ=Asia/Shanghai
WORKDIR /
ADD zone.tgz /
ADD localtime /etc/localtime
ADD yaml-v.tgz /
ENTRYPOINT ["/yaml-v"]
