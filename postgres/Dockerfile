FROM alpine:3.16

RUN apk add postgresql

COPY schema.sql /scaling/sql
COPY initialize.sh /scaling/sql
