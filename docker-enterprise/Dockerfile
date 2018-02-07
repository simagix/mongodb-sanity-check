FROM centos:centos7

MAINTAINER Ken Chen <simagix@gmail.com>

ADD ./mongodb-enterprise.repo /etc/yum.repos.d/

RUN yum -y update && \
    yum -y install mongodb-enterprise && \
    yum clean all && \
    mkdir -p /data/db && \
    chown -R mongod:mongod /data/db

VOLUME ["/data/db"]

EXPOSE 27017

CMD ["/usr/bin/mongod", "--dbpath", "/data/db"]
