// Package middleware define middleware propio para actuar de intermediario entre
// la llegada de una petición y su tratamiento final por un handler
package middleware

import (
	"backend/dao"
	"backend/globales"
	"backend/vo"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	NOMBRE_COOKIE_USUARIO             = "cookie_user" // El valor de una cookie de usuario consiste en "nombre_usuario'|'id'
	SEPARADOR_VALOR_COOKIE_USUARIO    = '|'
	LONGITUD_ID_COOKIE_USUARIO        = 128
	TIEMPO_EXPIRACION_COOKIES_USUARIO = 15 * 24 * time.Hour // 15 días
)

// Función que devuelve una función de middleware
func MiddlewareSesion() func(next http.Handler) http.Handler {
	// next es el handler (o middleware) siguiente a éste middleware
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Escribe directamente y luego deja escribir al handler a continuación,
			// en realidad se leerían cookies y se serviría contenido diferente o
			// dejaría pasar al handler, etc.
			cookie, cookieRequest, err := CargarCookieUsuario(r)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("Error al cargar cookie de usuario:", err)
			} else {
				if cookie.Expires.Before(time.Now()) && (cookieRequest.Raw == cookie.Raw) {
					// Cortamos la cadena en este punto, porque la cookie es inválida
					w.WriteHeader(http.StatusUnauthorized)
					log.Printf("Cookie expirada o inválida")
				} else {
					log.Println("Enhorabuena, estás logueado :)")
					// Deja pasar al handler
					next.ServeHTTP(w, r)
				}
			}
		}

		return http.HandlerFunc(fn)
	}
}

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

func CargarCookieUsuario(request *http.Request) (cookie http.Cookie, cookieRequest http.Cookie, err error) {
	for _, c := range request.Cookies() {
		if c.Name == NOMBRE_COOKIE_USUARIO { // Es una cookie de usuario
			// Obtener el usuario del valor de la cookie
			nombre := c.Value[:strings.IndexRune(c.Value, SEPARADOR_VALOR_COOKIE_USUARIO)]

			usuarioVO := vo.Usuario{0, "", nombre, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

			cookie, err = dao.ConsultarCookie(globales.Db, &usuarioVO)

			cookieRequest = *c

			break
		}
	}

	return cookie, cookieRequest, err
}

func GenerarCookieUsuario(writer *http.ResponseWriter, nombreUsuario string) (err error) {
	expiracion := time.Now().Add(TIEMPO_EXPIRACION_COOKIES_USUARIO)

	if !strings.Contains(nombreUsuario, "|") {
		valorCookie := nombreUsuario + string(SEPARADOR_VALOR_COOKIE_USUARIO) + RandStringRunes()

		cookie := http.Cookie{Name: NOMBRE_COOKIE_USUARIO, Value: valorCookie, Expires: expiracion}
		http.SetCookie(*writer, &cookie)

		usuarioVO := vo.Usuario{0, "", nombreUsuario, "", "", cookie, 0, 0, 0, 0, 0}

		err = dao.InsertarCookie(globales.Db, &usuarioVO)
	} else {
		log.Println(`Se ha proporcionado un nombre de usuario que contiene el carácter separador. 
                        No se ha generado ninguna cookie.`)
	}

	return err
}

// Genera un ID de una cookie aleatorio, de longitud LONGITUD_ID_COOKIE_USUARIO
func RandStringRunes() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	idCookie := make([]rune, LONGITUD_ID_COOKIE_USUARIO)
	for i := range idCookie {
		idCookie[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(idCookie)
}

func borrarCookieUsuario(writer *http.ResponseWriter, idUsuario string) {
	// db

	// Sobreescribe la cookie de usuario en la respuesta por la misma sin valor y expirando automáticamente
	cookie := http.Cookie{Name: NOMBRE_COOKIE_USUARIO, Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(*writer, &cookie)
}
