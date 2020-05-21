FROM edoburu/pgbouncer:1.12.0
MAINTAINER Failless
COPY ./configs/pgbouncer/pgbouncer.ini ./configs/pgbouncer/userlist.txt /etc/pgbouncer/

USER root
RUN mkdir -p /etc/pgbouncer /var/log/pgbouncer /var/run/pgbouncer && \
           chown -R postgres /var/run/pgbouncer /etc/pgbouncer /var/log/pgbouncer

ARG dbuser
ARG dbpasswd
ARG dbname

RUN echo "[databases]" >> /etc/pgbouncer/pgbouncer.ini && \
    echo "eventum = host=eventumdb port=5432 dbname=${dbname} user=${dbuser} password=${dbpasswd}" >> /etc/pgbouncer/pgbouncer.ini &&\
    echo "\"${dbuser}\" \"${dbpasswd}\"" >> /etc/pgbouncer/userlist.txt

USER postgres
EXPOSE 6432

CMD ["pgbouncer", "/etc/pgbouncer/pgbouncer.ini"]