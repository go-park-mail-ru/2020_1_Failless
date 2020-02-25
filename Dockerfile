FROM golang:1.12-stretch AS lang
WORKDIR /home/eventum
COPY . .
RUN go get -d && go build -v


FROM ubuntu:18.04
MAINTAINER Failless
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y gnupg
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

USER root
WORKDIR /home/eventum
RUN cd /home/eventum
COPY . .

RUN apt-get -y update
RUN apt-get -y install apt-transport-https git wget
RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main' >> /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update
ENV PGVERSION 11
RUN apt-get -y install postgresql-$PGVERSION postgresql-contrib

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER eventum WITH SUPERUSER PASSWORD 'eventum';" &&\
    createdb -O park park_forum && psql -d park_forum -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql eventum -a -f ./configs/migrations/init.sql &&\
    /etc/init.d/postgresql stop
EXPOSE 5432
EXPOSE 5000

WORKDIR /home/eventum
COPY --from=lang /home/eventum .

CMD service postgresql start && ./eventum
