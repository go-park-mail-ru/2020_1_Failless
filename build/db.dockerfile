FROM ubuntu:18.04
MAINTAINER Failless
ENV DEBIAN_FRONTEND noninteractive

USER root
WORKDIR /home/eventum
RUN cd /home/eventum
COPY ./scripts .

RUN apt-get -y update
RUN apt-get -y install apt-transport-https git wget
RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main' >> /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update

ENV PGVERSION 12
RUN apt-get -y install postgresql-$PGVERSION postgresql-contrib
RUN apt install -y postgis &&\
    apt install -y myspell-ru

ENV DICT /usr/share/hunspell/ru_RU
RUN iconv -f koi8-r -t utf-8 -o /usr/share/postgresql/$PGVERSION/tsearch_data/russian.affix $DICT.aff &&\
    iconv -f koi8-r -t utf-8 -o /usr/share/postgresql/$PGVERSION/tsearch_data/russian.dict  $DICT.dic

ARG dbuser
ARG dbpasswd
ARG dbname
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER ${dbuser} WITH SUPERUSER PASSWORD '${dbpasswd}';" &&\
    createdb -O $dbuser $dbname && psql -d eventum -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql eventum -a -f ./scripts/init.sql &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

CMD service postgresql start
