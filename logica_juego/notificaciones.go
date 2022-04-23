package logica_juego

import "encoding/gob"

// NotificacionAmistad representa una notificación de solicitud de amistad pendiente
//
// Ejemplo en JSON:
//    {
//        "IDNotificacion": 1,
//        "Jugador": 		"usuario6"
//    }
type NotificacionAmistad struct {
	IDNotificacion int    // 0
	Jugador        string // Jugador que ha enviado la solicitud de amistad
}

func NewNotificacionAmistad(jugador string) NotificacionAmistad {
	return NotificacionAmistad{IDNotificacion: NOTIFICACION_AMISTAD, Jugador: jugador}
}

// NotificacionTurno representa una notificación de que es el turno del jugador en
// la partida en la que está jugando
//
// Ejemplo en JSON:
//    {
//        "IDNotificacion":	1,
//        "JugadorPrevio":	"usuario6"
//    }
type NotificacionTurno struct {
	IDNotificacion int    // 1
	JugadorPrevio  string // Jugador del turno anterior
}

func NewNotificacionTurno(jugadorPrevio string) NotificacionTurno {
	return NotificacionTurno{IDNotificacion: NOTIFICACION_TURNO, JugadorPrevio: jugadorPrevio}
}

// NotificacionPuntosObtenidos representa una notificación de obtención de nuevos puntos, por ganar o perder una partida
//
// Ejemplo en JSON:
//    {
//        "IDNotificacion":	2,
//        "Puntos":	"usuario6",
//		  "PartidaGanada": false
//    }
type NotificacionPuntosObtenidos struct {
	IDNotificacion int  // 2
	Puntos         int  // Puntos obtenidos
	PartidaGanada  bool // Razón para la obtención de los puntos (false: por perder una partida, true: por ganar una partida)
}

func NewNotificacionPuntosObtenidos(puntos int, partidaGanada bool) NotificacionPuntosObtenidos {
	return NotificacionPuntosObtenidos{IDNotificacion: NOTIFICACION_PUNTOS, Puntos: puntos, PartidaGanada: partidaGanada}
}

// NotificacionExpulsion representa una notificación de que se ha sido expulsado de una partida por inactividad
//
// Ejemplo en JSON:
//    {
//        "IDNotificacion":	3,
//    }
type NotificacionExpulsion struct {
	IDNotificacion int // 3
}

func NewNotificacionExpulsion() NotificacionExpulsion {
	return NotificacionExpulsion{IDNotificacion: NOTIFICACION_EXPULSION}
}

func RegistrarNotificaciones() {
	gob.Register(NotificacionAmistad{})
	gob.Register(NotificacionTurno{})
	gob.Register(NotificacionPuntosObtenidos{})
	gob.Register(NotificacionExpulsion{})
}
