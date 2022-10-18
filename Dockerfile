FROM golang:latest

ENV TZ Asia/Shanghai

WORKDIR /weblog
ADD . /weblog

EXPOSE 8090
ENTRYPOINT ["./weblog"]