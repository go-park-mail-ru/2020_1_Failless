FROM edoburu/pgbouncer:1.12.0
COPY ./configs/pgbouncer/pgbouncer.ini ./configs/pgbouncer/userlist.txt /etc/pgbouncer/

USER root
RUN mkdir -p /etc/pgbouncer /var/log/pgbouncer /var/run/pgbouncer && \
           chown -R postgres /var/run/pgbouncer /etc/pgbouncer /var/log/pgbouncer
WORKDIR /home/eventum
RUN cd /home/eventum

USER postgres
EXPOSE 6432
VOLUME /etc/pgbouncer

CMD ["pgbouncer", "/etc/pgbouncer/pgbouncer.ini"]