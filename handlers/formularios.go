package handlers

import (
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Registro atiende respuestas de un formulario de campos 'nombre', 'email' y 'password'
// y registra un usuario acordemente. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
//
// Ruta: /registro
func Registro(writer http.ResponseWriter, request *http.Request) {
	nombre := request.FormValue("nombre")
	email := request.FormValue("email")
	password := request.FormValue("password")

	if dao.ExisteEmail(globales.Db, email) {
		devolverError(writer, errors.New("El email introducido ya está registrado"))
	}

	if dao.ExisteUsuario(globales.Db, nombre) {
		devolverError(writer, errors.New("El nombre de usuario introducido ya existe"))
	}

	hash, err := hashPassword(password)

	if err != nil {
		devolverError(writer, errors.New("Se ha producido un error al procesar los datos."))
	} else {
		usuarioVO := vo.Usuario{email, nombre, hash, "", http.Cookie{}, 0, 0, 0, 0, 0}
		err = dao.InsertarUsuario(globales.Db, &usuarioVO)

		if err != nil {
			devolverErrorSQL(writer)
		} else {
			err, cookie := middleware.GenerarCookieUsuario(&writer, nombre)
			if err != nil {
				devolverErrorSQL(writer)
			} else {
				writer.Write([]byte(cookie.String())) // La escribe en el body directamente
				escribirHeaderExito(writer)
			}
		}
	}
}

// Login atiende respuestas de un formulario de campos 'nombre' y 'password'
// y loguea a un usuario acordemente. Responde con status 200 y una cookie de usuario
// si ha habido éxito, o status 500 si ha habido un error junto a su motivo en el cuerpo.
//
// Ruta: /login
func Login(writer http.ResponseWriter, request *http.Request) {
	nombre := request.FormValue("nombre")
	password := request.FormValue("password")

	usuarioVO := vo.Usuario{"", nombre, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	hashDB, err := dao.ConsultarPasswordHash(globales.Db, &usuarioVO)

	existe := bcrypt.CompareHashAndPassword([]byte(hashDB), []byte(password))

	if err != nil || existe != nil {
		devolverError(writer, errors.New("Se ha producido un error al procesar los datos."))
	} else {
		err, cookie := middleware.GenerarCookieUsuario(&writer, nombre)
		if err != nil {
			devolverErrorSQL(writer)
		} else {
			writer.Write([]byte(cookie.String())) // La escribe en el body directamente
			escribirHeaderExito(writer)
		}
	}
}
