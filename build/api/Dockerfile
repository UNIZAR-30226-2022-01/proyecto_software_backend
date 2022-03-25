# syntax=docker/dockerfile:1
FROM scratch

# Copia todos los ficheros
COPY ./backend /     

# 80 para producción
EXPOSE 8090                 

ENTRYPOINT ["/backend", "-api"]
#CMD ["/backend"]
