FROM golang:1.13-stretch AS lang
WORKDIR /home/eventum
COPY . .
RUN go build -o bin/eventum

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
ENV PGVERSION 12
RUN apt-get -y install postgresql-$PGVERSION postgresql-contrib
RUN apt install -y postgis

USER postgres
ARG dbuser
ARG dbpasswd
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER ${dbuser} WITH SUPERUSER PASSWORD '${dbpasswd}';" &&\
    createdb -O $dbuser eventum && psql -d eventum -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql eventum -a -f ./configs/migrations/init.sql &&\
    /etc/init.d/postgresql stop
RUN apt install myspell-ru
RUN cd /usr/share/postgresql/$PGVERSION/tsearch_data
RUN ls
RUN DICT=/usr/share/hunspell/ru_RU
RUN iconv -f koi8-r -t utf-8 -o russian.affix $DICT.aff &&\
    iconv -f koi8-r -t utf-8 -o russian.dict  $DICT.dic

EXPOSE 5432
EXPOSE 5000

WORKDIR /home/eventum
COPY --from=lang /home/eventum .

CMD service postgresql start && ./bin/eventum
