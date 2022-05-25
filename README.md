# proyecto_software_backend
Repositorio para el backend de la asignatura de Proyecto Software


# Requisitos
- Compilador de go versión `>1.17`
- `docker` y `docker-compose`
- Para construir los servidores web, exclusivamente, `npm`

# Compilación y despliegue
El despliegue se encuentra automatizado y dividido entre despliegue real (rama ```main```) y despliegue de pruebas local (rama ```pruebas_local```)
- Por defecto, el servidor atenderá por el puerto ```443``` mediante `HTTPS`
- En el caso de pruebas en local, el servidor API atenderá mediante HTTP` por el puerto ```8090``` y el de Angular o React por ```8080```

## Servidor API
- Acceder a ```scripts/api```  y ejecutar el fichero ```crear_contenedores.sh```
- Se necesita tener todos los ficheros de variables de entorno en la carpeta y certificados en la carpeta ```envfiles```, que poseen los integrantes del grupo por privado
    -  ```mail.env```: Usuarios, contraseñas y nombres para el servicio de correos
    -   ```postgres.env```: Usuario y contraseña de la BD
    -   ```servidor.env```: Puertos y direcciones de la BD
    -   ```dns.env```: Nombres DNS para los que servir certificados TLS, en caso de necesitar usar ACME para obtener certificados
    -  ```clave_tls.key```: Clave del certificado TLS a usar
    -  ```cert_tls.key```: Certificado TLS a usar
  
 ## Servidor Angular / React
- Acceder a ```scripts/webserver_angular``` o ```scripts/webserver_react``` y ejecutar el fichero ```crear_contenedores.sh```
- Se necesita tener todos los ficheros de variables de entorno en la carpeta y certificados en la carpeta ```envfiles```, que poseen los integrantes del grupo por privado
    -   ```servidor.env```: Puertos y direcciones de la BD
    -   ```dns.env```: Nombres DNS para los que servir certificados TLS, en caso de necesitar usar ACME para obtener certificados
    -  ```clave_tls.key```: Clave del certificado TLS a usar
    -  ```cert_tls.key```: Certificado TLS a usar
