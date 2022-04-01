// Package middleware define middleware propio para actuar de intermediario entre
// la llegada de una petición y su tratamiento final por un handler, así como
// funciones auxiliares relacionadas con el mismo
package middleware

import (
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	NOMBRE_COOKIE_USUARIO             = "cookie_user" // Valor de una cookie de usuario: "nombre_usuario'|'id'
	SEPARADOR_VALOR_COOKIE_USUARIO    = '|'
	LONGITUD_ID_COOKIE_USUARIO        = 128
	TIEMPO_EXPIRACION_COOKIES_USUARIO = 15 * 24 * time.Hour // 15 días
)

// MiddlewareSesion devuelve un middleware que comprueba la existencia de una cookie de usuario válida antes de
// permitir a la URL dada, y deniega si no existe.
func MiddlewareSesion() func(next http.Handler) http.Handler {
	// next es el handler (o middleware) siguiente a éste middleware
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cookie, cookieRequest, err := CargarCookieUsuario(r)

			if err != nil {
				w.Write([]byte("Ha ocurrido un error al cargar la cookie de usuario, por favor, inicie sesión de nuevo."))
				w.WriteHeader(http.StatusUnauthorized)
				log.Println("Error al cargar cookie de usuario:", err)
			} else {
				if cookie.Expires.Before(time.Now()) && (cookieRequest.Raw == cookie.Raw) {
					w.Write([]byte("Tu sesión ha caducado, por favor, inicie sesión de nuevo."))
					// Se corta la cadena en este punto, porque la cookie es inválida
					w.WriteHeader(http.StatusUnauthorized)
					log.Println("Detectada cookie expirada o inválida:", cookie.Raw)
				} else {
					// Deja pasar al siguiente handler
					next.ServeHTTP(w, r)
				}
			}
		}

		return http.HandlerFunc(fn)
	}
}

// ObtenerUsuarioCookie devuelve el nombre de usuario almacenado en una cookie de usuario de la petición, si existe.
func ObtenerUsuarioCookie(request *http.Request) (nombre string) {
	for _, c := range request.Cookies() {
		if c.Name == NOMBRE_COOKIE_USUARIO { // Es una cookie de usuario
			// Obtener el usuario del valor de la cookie
			nombre = c.Value[:strings.IndexRune(c.Value, SEPARADOR_VALOR_COOKIE_USUARIO)]
			break
		}
	}

	return nombre
}

// CargarCookieUsuario devuelve la cookie de usuario almacenada en una cookie de usuario de la petición y la equivalente
// almacenada. Devuelve error en caso de no encontrarse alguna de las dos.
func CargarCookieUsuario(request *http.Request) (cookie http.Cookie, cookieRequest http.Cookie, err error) {
	encontrada := false

	for _, c := range request.Cookies() {
		if c.Name == NOMBRE_COOKIE_USUARIO { // Es una cookie de usuario
			// Obtener el usuario del valor de la cookie
			nombre := c.Value[:strings.IndexRune(c.Value, SEPARADOR_VALOR_COOKIE_USUARIO)]

			usuarioVO := vo.Usuario{"", nombre, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

			cookie, err = dao.ConsultarCookie(globales.Db, &usuarioVO)

			cookieRequest = *c

			encontrada = true
			break
		}
	}
	if !encontrada {
		// No se ha encontrado la cookie en la BD
		err = errors.New("No se ha encontrado la cookie en la BD")
	}

	return cookie, cookieRequest, err
}

// GenerarCookieUsuario genera una cookie para el nombre de usuario dado, la devuelve y la almacena. En caso de fallo o
// usuario no existente devuelve un error.
func GenerarCookieUsuario(writer *http.ResponseWriter, nombreUsuario string) (err error, cookie http.Cookie) {
	expiracion := time.Now().Add(TIEMPO_EXPIRACION_COOKIES_USUARIO)

	if !strings.Contains(nombreUsuario, "|") {
		valorCookie := nombreUsuario + string(SEPARADOR_VALOR_COOKIE_USUARIO) + RandStringRunes()

		cookie := http.Cookie{Name: NOMBRE_COOKIE_USUARIO, Value: valorCookie, Expires: expiracion}
		http.SetCookie(*writer, &cookie)
		// TODO: Por refinar una vez realizadas las pruebas de despliegue en un dominio real
		//cookie.SameSite = 4 // None
		//cookie.Domain = "localhost"
		usuarioVO := vo.Usuario{"", nombreUsuario, "", "", cookie, 0, 0, 0, 0, 0}

		err = dao.InsertarCookie(globales.Db, &usuarioVO)

		return err, cookie
	} else {
		// No debería ocurrir
		log.Println(`Se ha proporcionado un nombre de usuario que contiene el carácter separador. 
                        No se ha generado ninguna cookie.`)

		return err, cookie
	}

}

// BorrarCookieUsuario borrar la cookie para el nombre de usuario dado en el cliente y en el almacén para el nombre de
// usuario dado. En caso de fallo o usuario no existente devuelve un error.
func BorrarCookieUsuario(writer *http.ResponseWriter, nombreUsuario string) (err error) {
	// Sobreescribe la cookie de usuario en la respuesta por la misma sin valor y expirando automáticamente
	cookie := http.Cookie{Name: NOMBRE_COOKIE_USUARIO, Value: "", Expires: time.Unix(0, 0)}

	usuarioVO := vo.Usuario{"", nombreUsuario, "", "", cookie, 0, 0, 0, 0, 0}

	err = dao.InsertarCookie(globales.Db, &usuarioVO)

	http.SetCookie(*writer, &cookie)

	return err
}

// RandStringRunes genera un ID de una cookie aleatorio, de longitud LONGITUD_ID_COOKIE_USUARIO.
func RandStringRunes() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	idCookie := make([]rune, LONGITUD_ID_COOKIE_USUARIO)
	for i := range idCookie {
		idCookie[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(idCookie)
}
