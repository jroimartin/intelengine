FROM ubuntu:latest

MAINTAINER Roi Martin <jroi.martin@gmail.com>

ADD intelsrv.cfg /etc/
ADD intelsrv start-intelsrv.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/intelsrv /usr/local/bin/start-intelsrv.sh

RUN useradd -r -s /bin/bash intelsrv
USER intelsrv
ENTRYPOINT ["/usr/local/bin/start-intelsrv.sh"]

EXPOSE 8001
