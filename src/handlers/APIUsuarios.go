package handlers

import (
	"backend/dao"
	"backend/globales"
	"backend/middleware"
	"backend/vo"
	"encoding/json"
	"github.com/go-chi/chi/v5"
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
		devolverErrorSQL(writer)
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
		devolverErrorSQL(writer)
	}

	devolverExito(writer)
}

// RechazarSolicitudAmistad rechaza una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
func RechazarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")
	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.RechazarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverErrorSQL(writer)
	}

	devolverExito(writer)
}

// ListarAmigos devuelve una lista con los nombres de los amigos del usuario que genera la solicitud
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará códgo 200
// Dicha lista se devuelve en el siguiente formato JSON:
//	[ string, string, ...]
func ListarAmigos(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	usuario := vo.Usuario{NombreUsuario: nombreUsuario}
	amigos, err := dao.ObtenerAmigos(globales.Db, &usuario)
	if err != nil {
		devolverErrorSQL(writer)
	}

	var amigosString []string
	for _, a := range amigos {
		amigosString = append(amigosString, a.NombreUsuario)
	}

	writer.Header().Set("Content-Type", "application/json")
	envioAmigos := vo.ElementoListaNombresUsuario{Nombres: amigosString}
	err = json.NewEncoder(writer).Encode(envioAmigos)
}

// ObtenerPerfilUsuario devuelve la información del perfil de un usuario, definido como parte de la URL
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará códgo 200
// El formato de la respuesta JSON es el siguiente:
// [
//	"Email": string
//	"Nombre": string
//	"Biografia": string
// 	"PartidasGanadas": int
// 	"PartidasTotales": int
// 	"Puntos": int
// 	"ID_dado": int
// 	"ID_ficha": int
// ]
func ObtenerPerfilUsuario(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := chi.URLParam(request, "nombre")
	usuario, err := dao.ObtenerUsuario(globales.Db, nombreUsuario)
	if err != nil {
		devolverErrorSQL(writer)
	}

	envioUsuario := transformaAElementoListaUsuarios(usuario)
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(envioUsuario)
}

// ObtenerUsuariosSimilares devuelve una lista de nombres de usuario que coincidan con un patrón,
// especificado en al parámetro "patron" de la URL, ordenados alfabéticamente
// Los nombres de usuario coincidirán con dicho patrón o empezarán por él
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará códgo 200
// El formato de la respuesta JSON es el siguiente:
// [string, string, ...]
func ObtenerUsuariosSimilares(writer http.ResponseWriter, request *http.Request) {
	patron := chi.URLParam(request, "patron")
	usuarios, err := dao.ObtenerUsuariosSimilares(globales.Db, patron)
	if err != nil {
		devolverErrorSQL(writer)
	}

	envioUsuarios := vo.ElementoListaNombresUsuario{Nombres: usuarios}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(envioUsuarios)
}

func transformaAElementoListaUsuarios(usuario vo.Usuario) vo.ElementoListaUsuarios {
	return vo.ElementoListaUsuarios{
		NombreUsuario:   usuario.NombreUsuario,
		Email:           usuario.Email,
		Biografia:       usuario.Biografia,
		PartidasGanadas: usuario.PartidasGanadas,
		PartidasTotales: usuario.PartidasTotales,
		Puntos:          usuario.Puntos,
		ID_dado:         usuario.ID_dado,
		ID_ficha:        usuario.ID_ficha,
	}
}

// ObtenerNotificaciones devuelve un listado codificado en JSON de notificaciones
// a mostrar, relativas al usuario que lo solicita.
func ObtenerNotificaciones(writer http.ResponseWriter, request *http.Request) {
	// TODO
}
