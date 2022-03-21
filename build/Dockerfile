# syntax=docker/dockerfile:1
FROM scratch
#FROM golang:alpine
# Directorio de trabajo para RUN, CMD, COPY, etc.
#WORKDIR backend  

# Copia todos los ficheros
COPY ./backend /
COPY ./web /web          

# 80 para producción
EXPOSE 8080                 

CMD [ "/backend" ]
