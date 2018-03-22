FROM ubuntu

MAINTAINER WSL wu12490@gmail.com

RUN apt-get update
RUN apt-get install -y dnsmasq

COPY ./config/* /etc/
COPY ./dns-driver /usr/bin/

EXPOSE 8082
EXPOSE 53

ENTRYPOINT sh /etc/start.sh
