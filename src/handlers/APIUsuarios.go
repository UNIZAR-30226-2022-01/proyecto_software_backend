package handlers

import (
	"backend/dao"
	"backend/globales"
	"backend/middleware"
	"backend/vo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

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
