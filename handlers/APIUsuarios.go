package handlers

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"github.com/go-chi/chi/v5"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

	if nombreUsuarioEmisor == nombreUsuarioReceptor {
		devolverError(writer, errors.New("No puedes enviarte una solicitud de amistad a ti mismo"))
		return
	}

	usuarioEmisor := vo.Usuario{"", nombreUsuarioEmisor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	usuarioReceptor := vo.Usuario{"", nombreUsuarioReceptor, "", "", http.Cookie{}, 0, 0, 0, 0, 0}

	err := dao.CrearSolicitudAmistad(globales.Db, &usuarioEmisor, &usuarioReceptor)
	if err != nil {
		devolverErrorSQL(writer)
		return
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
		return
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
		return
	}

	escribirHeaderExito(writer)
}

// EliminarAmigo elimina una relación de amistad entre dos usuarios, el nombre del usuario a borrar de tu lista
// se especificará como parte de la URL. Devolverá status 500 en caso de error, 200 en cualquier otro caso.
//
// Ruta: /api/eliminarAmigo/{nombre}
// Tipo: GET
func EliminarAmigo(writer http.ResponseWriter, request *http.Request) {
	nombreAmigo := chi.URLParam(request, "nombre")
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	usuario1 := vo.Usuario{NombreUsuario: nombreUsuario}
	usuario2 := vo.Usuario{NombreUsuario: nombreAmigo}

	err := dao.EliminarAmigo(globales.Db, &usuario1, &usuario2)
	if err != nil {
		devolverErrorSQL(writer)
		return
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
		return
	}

	var listaAmigos []string
	for _, a := range amigos {
		listaAmigos = append(listaAmigos, a.NombreUsuario)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(listaAmigos)
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
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(pendientes)
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
// 	   "ID_avatar": int
//	   "EsAmigo": bool
//	   "SolicitudRecibida": bool
//	   "SolicitudPendiente": bool
//    }
//
// Ruta: /api/obtenerPerfil/{nombre}
// Tipo: GET
func ObtenerPerfilUsuario(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := chi.URLParam(request, "nombre")
	usuario, err := dao.ObtenerUsuario(globales.Db, nombreUsuario)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	envioUsuario := transformaAElementoListaUsuarios(usuario)
	// Se comprueba si es amigo del usuario solicitante o no
	usuarioSolicitante := middleware.ObtenerUsuarioCookie(request)
	amigos, err := dao.ObtenerAmigos(globales.Db, &vo.Usuario{NombreUsuario: usuarioSolicitante})
	for _, amigo := range amigos {
		if amigo.NombreUsuario == nombreUsuario {
			envioUsuario.EsAmigo = true
			break
		}
	}

	// Si no son amigos, comprobamos si la solicitud de amistad está pendiente
	if !envioUsuario.EsAmigo {
		envioUsuario.SolicitudRecibida, envioUsuario.SolicitudPendiente =
			comprobarEstadoSolicitud(usuarioSolicitante, nombreUsuario)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(envioUsuario)
	escribirHeaderExito(writer)
}

// ObtenerUsuariosSimilares devuelve una lista de nombres de usuario que coincidan con un patrón,
// especificado en al parámetro "patron" de la URL, ordenados alfabéticamente e indicando si son
// amigos del usuario que lo solicita o no
// Los nombres de usuario coincidirán con dicho patrón o empezarán por él
// Si ocurre algún error durante el procesamiento, enviará código de error 500
// En cualquier otro caso, enviará código 200
// El formato de la respuesta JSON es el siguiente:
//    [
//        {
//            "Nombre": string,
//            "EsAmigo": bool,
//			  "SolicitudRecibida": bool,
//		 	  "SolicitudPendiente": bool
//        },
//        {
//            "Nombre": string,
//            "EsAmigo": bool,
//			  "SolicitudRecibida": bool,
//		 	  "SolicitudPendiente": bool
//        },
//		  ...
//    ]
//
// Ruta: /api/obtenerUsuariosSimilares/{patron}
// Tipo: GET
func ObtenerUsuariosSimilares(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)

	patron := chi.URLParam(request, "patron")
	usuarios, err := dao.ObtenerUsuariosSimilares(globales.Db, patron)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	var envioUsuarios []vo.ElementoListaUsuariosSimilares

	// TODO: Comprobar eficiencia con tests de carga, es O(n^2), llevar a DB si no
	amigos, err := dao.ObtenerAmigos(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})
	for _, usuario := range usuarios {
		esAmigo := false
		for i, amigo := range amigos {
			if usuario == amigo.NombreUsuario {
				esAmigo = true
				amigos = append(amigos[:i], amigos[i+1:]...) // Lo elimina de la lista, no hay que comprobarlo de nuevo
			}
		}

		// Si no son amigos, comprobamos solicitudes de amistad pendientes
		solicitudRecibida := false
		solicitudPendiente := false
		if !esAmigo {
			solicitudRecibida, solicitudPendiente = comprobarEstadoSolicitud(nombreUsuario, usuario)
		}

		envioUsuarios = append(envioUsuarios, vo.ElementoListaUsuariosSimilares{Nombre: usuario, EsAmigo: esAmigo,
			SolicitudRecibida: solicitudRecibida, SolicitudPendiente: solicitudPendiente})
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(envioUsuarios)
	escribirHeaderExito(writer)
}

// ObtenerRanking devuelve una lista de todos los usuarios registrados, ordenados por el número de partidas ganadas.
// Dicha lista se devolverá en el siguiente formato JSON:
// [
//  {
//   NombreUsuario: string
//   PartidasGanadas: int
//   PartidasTotales: int
// 	}, {...}, ...
// ]
//
// Devolverá status 500 en caso de error y status 200 en cualquier otro caso
//
// Ruta: /api/ranking
// Tipo: GET
func ObtenerRanking(writer http.ResponseWriter, request *http.Request) {
	ranking, err := dao.Ranking(globales.Db)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(ranking)
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
// Ruta: /api/obtenerNotificaciones
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

	// Añade las notificaciones con estado almacenadas en la base de datos
	err, notificacionesConEstado := dao.ObtenerNotificacionesConEstado(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario}, true)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	if len(notificacionesConEstado) > 0 {
		notificaciones = append(notificaciones, notificacionesConEstado)
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(notificaciones)
	escribirHeaderExito(writer)
}

// ObtenerNumeroNotificaciones devuelve el número de notificaciones que se van a recibir en
// ObtenerNotificaciones, sin consumirlas.
//
// El formato de la respuesta JSON es el siguiente:
//	int
//
// Ruta: /api/obtenerNumeroNotificaciones
// Tipo: GET
func ObtenerNumeroNotificaciones(writer http.ResponseWriter, request *http.Request) {
	contador := 0

	nombreUsuario := middleware.ObtenerUsuarioCookie(request)

	usuariosPendientes, err := dao.ConsultarSolicitudesPendientes(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})
	contador += len(usuariosPendientes)

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
			contador += 1
		}
	}

	// Añade las notificaciones con estado almacenadas en la base de datos
	err, notificacionesConEstado := dao.ObtenerNotificacionesConEstado(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario}, false)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	contador += len(notificacionesConEstado)

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(contador)
	escribirHeaderExito(writer)
}

// ResetearContraseña resetea la contraseña del usuario que tiene el token indicado.
// Devuelve status 500 en caso de error y status 200 en cualquier otro caso
//
// Entrada: formulario con campos "password" y "token"
//
// Ruta: /resetearPassword
// Tipo: POST
func ResetearContraseña(writer http.ResponseWriter, request *http.Request) {
	contraseña := request.FormValue("password")
	token := request.FormValue("token")

	hash, err := hashPassword(contraseña)

	err, usuario, expiracion := dao.ObtenerNombreExpiracionTokenResetPassword(globales.Db, token)
	if err != nil {
		devolverError(writer, errors.New("El token que has indicado no existe, ha expirado, o ha ocurrido un error inesperado. Por favor, inténtalo de nuevo más tarde."))
		return
	}

	err = dao.ResetearContraseña(globales.Db, usuario, hash)
	if err != nil || expiracion.Before(time.Now()) {
		devolverError(writer, errors.New("El token que has indicado no existe, ha expirado, o ha ocurrido un error inesperado. Por favor, inténtalo de nuevo más tarde."))
		return
	}

	escribirHeaderExito(writer)
}

// ObtenerTokenResetPassword envía por email al usuario indicado un link de reset de contraseña, si existe.
// Devuelve status 500 en caso de error y status 200 en cualquier otro caso
//
// Entrada: formulario con campo "usuario"
//
// Ruta: /obtenerTokenResetPassword
// Tipo: POST
func ObtenerTokenResetPassword(writer http.ResponseWriter, request *http.Request) {
	usuario := request.FormValue("usuario")

	err, token := dao.CrearTokenResetPassword(globales.Db, usuario)
	if err != nil {
		devolverError(writer, errors.New("El nombre de usuario que has indicado no existe, o ha ocurrido un error inesperado. Por favor, inténtalo de nuevo más tarde."))
		return
	}

	err, email := dao.ObtenerEmailUsuario(globales.Db, usuario)
	if err != nil {
		devolverError(writer, errors.New("El nombre de usuario que has indicado no existe, o ha ocurrido un error inesperado. Por favor, inténtalo de nuevo más tarde."))
	} else {
		m := gomail.NewMessage()
		m.SetHeader("From", os.Getenv(globales.DIRECCION_ENVIO_EMAILS))
		m.SetHeader("To", email)
		m.SetHeader("Subject", "¡"+usuario+", resetea tu contraseña aquí!")
		//m.SetBody("text/html", "<a>"+os.Getenv(globales.NOMBRE_DNS_API)+"/resetearPassword/</a>"+token)
		m.SetBody("text/html", "Has recibido este email porque has solicitado cambiar tu contraseña."+
			"La próxima vez que inicies sesión, indique que desea resetear la contraseña e introduzca el siguiente token en el campo en el que se lo indique: "+token)

		puerto, _ := strconv.Atoi(os.Getenv(globales.PUERTO_SMTP))
		d := gomail.NewDialer(os.Getenv(globales.HOST_SMTP), puerto, os.Getenv(globales.USUARIO_SMTP), os.Getenv(globales.PASS_SMTP))
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Evita problemas derivados de no tener certificados en el contenedor de Docker
		// Como no se puede verificar que el destino no existe, los errores al enviar correos se ignoran silenciosamente
		if err := d.DialAndSend(m); err != nil {
			log.Println("Error al enviar email de reset de contraseña a", usuario, ":", err)
		}

		escribirHeaderExito(writer)
	}
}

// ModificarBiografia permite al usuario modificar su biografía, especificando una nueva en el campo "biografia" del
// formulario enviado.
//
// Devuelve status 500 en caso de error y 200 en caso contrario
//
// Ruta: /api/modificarBiografia
// Tipo: Post
func ModificarBiografia(writer http.ResponseWriter, request *http.Request) {
	usuario := middleware.ObtenerUsuarioCookie(request)
	biografia := request.FormValue("biografia")

	err := dao.ModificarBiografia(globales.Db, &vo.Usuario{NombreUsuario: usuario}, biografia)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	escribirHeaderExito(writer)
}

// ModificarAspecto permite al usuario equipar un aspecto que haya comprado previamente. Para ello, especificará
// el aspecto con el identificador del objeto en la URL. En caso de que ese objeto no exista, o no esté en la colección
//del usuario, no lo podrá equipar.
//
// Devuelve status 500 en caso de error, 200 en cualquier otro caso
//
// Ruta: /api/modificarAspecto/{id_aspecto}
// Tipo: POST
func ModificarAspecto(writer http.ResponseWriter, request *http.Request) {
	usuario := middleware.ObtenerUsuarioCookie(request)
	idAspecto, err := strconv.Atoi(chi.URLParam(request, "id_aspecto"))
	if err != nil {
		devolverError(writer, errors.New("El identificador del aspecto debe ser un número natural"))
		return
	}

	aspecto, err := dao.ObtenerObjeto(globales.Db, idAspecto)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Comprobamos que el usuario tenga el aspecto o que es algún aspecto inicial
	if aspecto.Id > 2 {
		existe, err := dao.TieneObjeto(globales.Db, &vo.Usuario{NombreUsuario: usuario}, aspecto)
		if err != nil {
			devolverErrorSQL(writer)
			return
		}

		if !existe {
			devolverError(writer, errors.New("No puedes equipar un aspecto que no has comprado"))
			return
		}
	}

	switch aspecto.Tipo {
	case "avatar":
		// Modificamos el aspecto de las fichas del jugador
		err = dao.ModificarAvatar(globales.Db, &vo.Usuario{NombreUsuario: usuario}, aspecto)
	case "dado":
		// Modificamos el aspecto de las dados del jugador
		err = dao.ModificarDados(globales.Db, &vo.Usuario{NombreUsuario: usuario}, aspecto)
	default:
		devolverError(writer, errors.New("Aspecto no reconocido"))
		return
	}

	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	escribirHeaderExito(writer)
}
