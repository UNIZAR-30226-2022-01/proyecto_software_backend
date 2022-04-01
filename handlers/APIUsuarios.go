package handlers

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// EnviarSolicitudAmistad envía una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
//
// Ruta: /api/enviarSolicitudAmistad/{nombre}
// Tipo: POST
func EnviarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")
	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.CrearSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverErrorSQL(writer)
	}

	escribirHeaderExito(writer)
}

// AceptarSolicitudAmistad acepta una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
//
// Ruta: /api/aceptarSolicitudAmistad/{nombre}
// Tipo: POST
func AceptarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")
	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.AceptarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverErrorSQL(writer)
	}

	escribirHeaderExito(writer)
}

// RechazarSolicitudAmistad rechaza una solicitud de amistad entre el usuario que genera
// la solicitud y el indicado en el nombre. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
//
// Ruta: /api/rechazarSolicitudAmistad/{nombre}
// Tipo: POST
func RechazarSolicitudAmistad(writer http.ResponseWriter, request *http.Request) {
	nombreUsuarioReceptor := chi.URLParam(request, "nombre")
	nombreUsuarioEmisor := middleware.ObtenerUsuarioCookie(request)

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.RechazarSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverErrorSQL(writer)
	}

	escribirHeaderExito(writer)
}

// ListarAmigos devuelve una lista con los nombres de los amigos del usuario que genera la solicitud
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará códgo 200
// Dicha lista se devuelve en el siguiente formato JSON:
//	[ string, string, ...]
//
// Ruta: /api/listarAmigos
// Tipo: GET
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
	escribirHeaderExito(writer)
}

// ObtenerSolicitudesPendientes devuelve una lista de nombres de usuario a los
// que se les ha enviado una solicitud de amistad aún pendiente por aceptar o
// rechazar, codificada en JSON.
//
// El formato de la respuesta JSON es el siguiente:
// ["nombre1", "nombre2", ...]
//
// Ruta: /api/obtenerSolicitudesPendientes
// Tipo: GET
func ObtenerSolicitudesPendientes(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	usuario := vo.Usuario{NombreUsuario: nombreUsuario}
	pendientes, err := dao.ConsultarSolicitudesPendientes(globales.Db, &usuario)
	if err != nil {
		devolverErrorSQL(writer)
	}

	writer.Header().Set("Content-Type", "application/json")
	envioPendientes := vo.ElementoListaNombresUsuario{Nombres: pendientes}
	err = json.NewEncoder(writer).Encode(envioPendientes)
	escribirHeaderExito(writer)
}

// ObtenerPerfilUsuario devuelve la información del perfil de un usuario, definido como parte de la URL
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará códgo 200
//
// El formato de la respuesta JSON es el siguiente:
//    {
//	   "Email": string
//	   "Nombre": string
//	   "Biografia": string
// 	   "PartidasGanadas": int
// 	   "PartidasTotales": int
// 	   "Puntos": int
// 	   "ID_dado": int
// 	   "ID_ficha": int
//    }
//
// Ruta: /api/obtenerPerfil/{nombre}
// Tipo: GET
func ObtenerPerfilUsuario(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := chi.URLParam(request, "nombre")
	usuario, err := dao.ObtenerUsuario(globales.Db, nombreUsuario)
	if err != nil {
		devolverErrorSQL(writer)
	}

	envioUsuario := transformaAElementoListaUsuarios(usuario)
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(envioUsuario)
	escribirHeaderExito(writer)
}

// ObtenerUsuariosSimilares devuelve una lista de nombres de usuario que coincidan con un patrón,
// especificado en al parámetro "patron" de la URL, ordenados alfabéticamente
// Los nombres de usuario coincidirán con dicho patrón o empezarán por él
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará código 200
// El formato de la respuesta JSON es el siguiente:
//    [string, string, ...]
//
// Ruta: /api/obtenerUsuariosSimilares/{patron}
// Tipo: GET
func ObtenerUsuariosSimilares(writer http.ResponseWriter, request *http.Request) {
	patron := chi.URLParam(request, "patron")
	usuarios, err := dao.ObtenerUsuariosSimilares(globales.Db, patron)
	if err != nil {
		devolverErrorSQL(writer)
	}

	envioUsuarios := vo.ElementoListaNombresUsuario{Nombres: usuarios}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(envioUsuarios)
	escribirHeaderExito(writer)
}

// ObtenerNotificaciones devuelve un listado codificado en JSON de notificaciones
// a mostrar, relativas al usuario que lo solicita.
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará código 200 y la lista de notificaciones.
//
// El formato de la respuesta JSON es el siguiente:
//    [notificacion1..., notificacion2...]
//
// Ejemplo:
//    [{"IDNotificacion":0,"Jugador":"usuario2"}, {"IDNotificacion":1,"JugadorPrevio":"usuario6"}]
//
// La lista de notificaciones y su formato en JSON están disponibles en el módulo de logica_juego, en notificaciones.go
//
// Ruta: /api/obtenerNotificaciones/
// Tipo: GET
func ObtenerNotificaciones(writer http.ResponseWriter, request *http.Request) {
	var notificaciones []interface{}

	nombreUsuario := middleware.ObtenerUsuarioCookie(request)

	usuariosPendientes, err := dao.ConsultarSolicitudesPendientes(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})

	for _, usuario := range usuariosPendientes {
		notificaciones = append(notificaciones, logica_juego.NewNotificacionAmistad(usuario))
	}

	enPartida, err := dao.UsuarioEnPartida(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	if enPartida {
		idPartida, err := dao.PartidaUsuario(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})
		if err != nil {
			devolverErrorSQL(writer)
			return
		}

		partida, _ := globales.CachePartidas.ObtenerPartida(idPartida)
		if partida.Estado.Jugadores[partida.Estado.TurnoJugador] == nombreUsuario {
			turnoPrevio := partida.Estado.TurnoJugador - 1
			if turnoPrevio == -1 {
				turnoPrevio = len(partida.Estado.Jugadores) - 1
			}

			notificaciones = append(notificaciones, logica_juego.NewNotificacionTurno(partida.Estado.Jugadores[turnoPrevio]))
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(notificaciones)
	escribirHeaderExito(writer)
}
