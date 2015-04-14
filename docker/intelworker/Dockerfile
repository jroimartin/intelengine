FROM ubuntu:latest

MAINTAINER Roi Martin <jroi.martin@gmail.com>

ADD intelworker.cfg /etc/
ADD intelworker start-intelworker.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/intelworker /usr/local/bin/start-intelworker.sh

RUN useradd -r -s /bin/bash intelworker

RUN mkdir -p /commands
RUN chown intelworker:intelworker /commands
VOLUME ["/commands"]

USER intelworker
ENTRYPOINT ["/usr/local/bin/start-intelworker.sh"]
