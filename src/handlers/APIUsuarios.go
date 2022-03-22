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

// EnviarSolicitudAmistad envía una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
func EnviarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")
	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.CrearSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverError(writer, "EnviarSolicitudAmistad", err)
	}

	devolverExito(writer)
}

// AceptarSolicitudAmistad acepta una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
func AceptarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")
	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.AceptarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverError(writer, "AceptarSolicitudAmistad", err)
	}

	devolverExito(writer)
}

// RechazarSolicitudAmistad rechaza una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
func RechazarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")

	log.Println("Parámetros RechazarSolicitudAmistad:", nombreUsuarioReceptor)

	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.RechazarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverError(writer, "RechazarSolicitudAmistad", err)
	}

	devolverExito(writer)
}

// ObtenerNotificaciones devuelve un listado codificado en JSON de notificaciones
// a mostrar, relativas al usuario que lo solicita.
func ObtenerNotificaciones(writer http.ResponseWriter, request *http.Request) {
	// TODO
}
