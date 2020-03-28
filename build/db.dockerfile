FROM ubuntu:18.04
MAINTAINER Failless
ENV DEBIAN_FRONTEND noninteractive

USER root
WORKDIR /home/eventum
RUN cd /home/eventum
COPY ./scripts .

RUN apt-get -y update && apt-get install -y gnupg sed
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
    export user=$(echo ${dbuser} | sed -e 's|["'\'']||g'); psql --command "CREATE USER $user WITH SUPERUSER PASSWORD '${dbpasswd}';" &&\
    createdb -O $user $dbname && psql -d $dbname -c "CREATE EXTENSION IF NOT EXISTS citext;" &&\
    psql $dbname -a -f ./init.sql &&\
    /etc/init.d/postgresql stop

USER root
RUN echo "local all all md5" > /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "listen_addresses = '*'" >> /etc/postgresql/$PGVERSION/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432
USER postgres
CMD ["/usr/lib/postgresql/12/bin/postgres", "-D", "/var/lib/postgresql/12/main", "-c", "config_file=/etc/postgresql/12/main/postgresql.conf"]
