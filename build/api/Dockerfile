# syntax=docker/dockerfile:1
FROM golang:alpine

# Copia todos los ficheros
COPY ./backend /go
COPY ./assets /go/assets

EXPOSE 443              

ENTRYPOINT ["/go/backend", "-api"]
#CMD ["/backend"]
