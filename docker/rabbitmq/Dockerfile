FROM ubuntu:latest

MAINTAINER Roi Martin <jroi.martin@gmail.com>

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates wget
RUN echo 'deb http://www.rabbitmq.com/debian/ testing main' > /etc/apt/sources.list.d/10rabbitmq.list
RUN wget -O - 'https://www.rabbitmq.com/rabbitmq-signing-key-public.asc' | apt-key add -
RUN apt-get update && apt-get -y install rabbitmq-server
RUN rabbitmq-plugins enable rabbitmq_management && \
	echo '[{rabbit, [{loopback_users, []}]}].' > /etc/rabbitmq/rabbitmq.config

RUN mkdir -p /data/mnesia /data/log
RUN chown rabbitmq:rabbitmq /data/mnesia /data/log
ENV RABBITMQ_MNESIA_BASE /data/mnesia
ENV RABBITMQ_LOG_BASE /data/log
VOLUME ["/data/mnesia", "/data/log"]

ADD start-rabbitmq.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/start-rabbitmq.sh
ENTRYPOINT ["/bin/sh", "-c", "/usr/local/bin/start-rabbitmq.sh"]

EXPOSE 5672
EXPOSE 15672
