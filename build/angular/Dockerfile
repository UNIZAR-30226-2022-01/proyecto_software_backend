# syntax=docker/dockerfile:1
FROM golang:alpine

# Copia todos los ficheros
COPY ./backend /go
COPY ./web /go/web          

EXPOSE 443                 

ENTRYPOINT ["/backend", "-web"]
#CMD ["/backend"]
