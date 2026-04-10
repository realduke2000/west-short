FROM mirror.ccs.tencentyun.com/library/ubuntu:22.04
COPY repo/shortsrv /usr/local/bin/shortsrv
COPY repo/wshort.conf /etc/wshort.conf
EXPOSE 80
ENTRYPOINT ["/usr/local/bin/shortsrv"]