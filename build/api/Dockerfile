# syntax=docker/dockerfile:1
FROM scratch

# Copia todos los ficheros
COPY ./backend /
COPY ./web /web          

# 80 para producción
EXPOSE 8080                 

ENTRYPOINT ["/backend", "-api"]
#CMD ["/backend"]
