FROM postgres

COPY database/sql /sql
COPY database/allin.sh /
RUN chmod +x /allin.sh

RUN /allin.sh /docker-entrypoint-initdb.d/all.sql
RUN chmod +x /docker-entrypoint-initdb.d/*.sql
