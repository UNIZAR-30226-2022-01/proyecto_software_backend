// Package generales define handlers de páginas accesibles para cualquier usuario
package generales

import (
	"backend/dao"
	"backend/globales"
	"backend/middleware"
	"backend/vo"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	//"html/template"
	//"log"
	//"net/http"
	//"webserver/globales"
	//"webserver/vo"
)

// Struct de ejemplo a serializar a JSON
type Ejemplo struct {
	Campo1 string
	Campo2 int
}

func Registro(writer http.ResponseWriter, request *http.Request) {
	nombre := request.FormValue("nombre")
	email := request.FormValue("email")
	password := request.FormValue("password")

	log.Println("Parámetros:", nombre, email, password)

	hash, err := hashPassword(password)

	if err != nil {
		log.Println("Error al generar hash de contraseña para", nombre)
	} else {
		usuarioVO := vo.Usuario{0, email, nombre, hash, "", http.Cookie{}, 0, 0, 0, 0, 0}
		err = dao.InsertarUsuario(globales.Db, &usuarioVO)

		// TODO: Forma de comunicar error en registro
		if err != nil {
			log.Println("Error en registro:", err)
			http.Redirect(writer, request, "/menuRegistro", 307) // Código de redirección temporal
		} else {
			middleware.GenerarCookieUsuario(&writer, nombre)
		}
	}

	tmpl := template.Must(template.ParseFiles("web/index.html"))
	tmpl.Execute(writer, nil)
}

func Login(writer http.ResponseWriter, request *http.Request) {
	nombre := request.FormValue("nombre")
	password := request.FormValue("password")

	log.Println("Parámetros:", nombre, password)

	usuarioVO := vo.Usuario{0, "", nombre, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	hashDB, err := dao.ConsultarPasswordHash(globales.Db, &usuarioVO)

	existe := bcrypt.CompareHashAndPassword([]byte(hashDB), []byte(password))

	if err != nil || existe != nil {
		// TODO: Forma de comunicar error en login
		log.Println("Contraseña incorrecta")
		http.Redirect(writer, request, "/menuLogin", 307) // Código de redirección temporal
	} else {
		middleware.GenerarCookieUsuario(&writer, nombre)
	}

	tmpl := template.Must(template.ParseFiles("web/index.html"))
	tmpl.Execute(writer, nil)
}

func EnviarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")

	log.Println("Parámetros EnviarSolicitudAmistad:", nombreUsuarioReceptor)

	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{0, "", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{0, "", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.CrearSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		log.Println("Error en EnviarSolicitudAmistad", err)
	}

	// TODO: respuesta
}

func AceptarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")

	log.Println("Parámetros AceptarSolicitudAmistad:", nombreUsuarioReceptor)

	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{0, "", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{0, "", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.AceptarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		log.Println("Error en AceptarSolicitudAmistad", err)
	}

	// TODO: respuesta
}

func RechazarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")

	log.Println("Parámetros RechazarSolicitudAmistad:", nombreUsuarioReceptor)

	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{0, "", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{0, "", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.RechazarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		log.Println("Error en RechazarSolicitudAmistad", err)
	}

	// TODO: respuesta
}

func ServirJSON(writer http.ResponseWriter, request *http.Request) {
	ejemplo := Ejemplo{"campo1", 2}

	log.Println(json.Marshal(ejemplo))

	// Establece el contenido a servir como JSON y lo escribe
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(ejemplo)
}

func ServirImagen(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "unizzard.png")
}

// hashPassword crea un hash de clave utilizando bcrypt
// https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Coste fijo generoso
	return string(bytes), err
}
