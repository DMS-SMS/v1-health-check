FROM alpine
MAINTAINER Park, Jinhong <jinhong0719@naver.com>

COPY ./health-check ./health-check
ENTRYPOINT [ "/health-check" ]
