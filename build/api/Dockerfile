# syntax=docker/dockerfile:1
FROM golang:alpine

# Copia todos los ficheros
COPY ./backend /     
COPY ./assets /assets

EXPOSE 443              

ENTRYPOINT ["/backend", "-api"]
#CMD ["/backend"]
