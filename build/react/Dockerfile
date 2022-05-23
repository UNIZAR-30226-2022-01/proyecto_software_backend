# syntax=docker/dockerfile:1
FROM golang:alpine

# Copia todos los ficheros
COPY ./backend /go
COPY ./web /go/web
COPY ./cert_tls.pem /go/
COPY ./clave_tls.key /go/

EXPOSE 443                

ENTRYPOINT ["/go/backend", "-web"]
#CMD ["/backend"]
