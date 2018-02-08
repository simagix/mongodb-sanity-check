FROM centos:centos7

MAINTAINER Ken Chen <simagix@gmail.com>

ADD ./mongodb-enterprise.repo /etc/yum.repos.d/
ADD ./*.keyfile /etc/ssl/

RUN yum -y update && \
    yum -y install mongodb-enterprise && \
    yum clean all && \
    mkdir -p /data/db /etc/ssl && \
    chown -R mongod:mongod /data/db /etc/ssl && \
    chmod 600 /etc/ssl/*.keyfile

ADD ./mongod.conf /etc/

VOLUME ["/data/db"]

EXPOSE 27017

CMD ["/usr/bin/mongod", "-f", "/etc/mongod.conf"]
