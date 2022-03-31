package logica_juego

import "encoding/gob"

// NotificacionAmistad representa una notificación de solicitud de amistad pendiente
type NotificacionAmistad struct {
	IDNotificacion int    // 0
	Jugador        string // Jugador que ha enviado la solicitud de amistad
}

func NewNotificacionAmistad(jugador string) NotificacionAmistad {
	return NotificacionAmistad{IDNotificacion: NOTIFICACION_AMISTAD, Jugador: jugador}
}

// NotificacionTurno representa una notificación de que es el turno del jugador en
// la partida en la que está jugando
type NotificacionTurno struct {
	IDNotificacion int    // 1
	JugadorPrevio  string // Jugador del turno anterior
}

func NewNotificacionTurno(jugadorPrevio string) NotificacionTurno {
	return NotificacionTurno{IDNotificacion: NOTIFICACION_TURNO, JugadorPrevio: jugadorPrevio}
}

func RegistrarNotificaciones() {
	gob.Register(NotificacionAmistad{})
	gob.Register(NotificacionTurno{})
}
