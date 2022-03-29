package vo

import (
	//"backend/logica_juego"
	"net/http"
)

// Usuario es un objeto de usuario equivalente al del modelo de base de datos.
type Usuario struct {
	Email           string
	NombreUsuario   string
	PasswordHash    string
	Biografia       string
	CookieSesion    http.Cookie
	Puntos          int
	PartidasGanadas int
	PartidasTotales int
	ID_dado         int
	ID_ficha        int
}

type Partida struct {
	IdPartida          int
	EsPublica          bool
	PasswordHash       string
	EnCurso            bool
	MaxNumeroJugadores int

	// TODO: representar chat de la partida
	Mensajes []Mensaje
	Estado   EstadoPartida
}

type Mensaje struct {
	// TODO
}

// CrearPartida devuelve una partida sin estado ni ID asignado.
func CrearPartida(esPublica bool, passwordHash string, maxNumeroJugadores int) *Partida {
	partida := Partida{
		IdPartida:          0,
		EsPublica:          esPublica,
		PasswordHash:       passwordHash,
		EnCurso:            false,
		MaxNumeroJugadores: maxNumeroJugadores,
	}

	return &partida
}

// IniciarPartida marca una partida como iniciada y crea un estado para ella con los jugadores indicados
func (p *Partida) IniciarPartida(jugadores []string) {
	p.EnCurso = true
	p.Estado = *CrearEstadoPartida(jugadores, p.MaxNumeroJugadores)
}

///////////////////////////////////////////
// Structs de respuesta a serializar a JSON
///////////////////////////////////////////

type ElementoListaPartidas struct {
	IdPartida          int
	EsPublica          bool
	NumeroJugadores    int
	MaxNumeroJugadores int
	AmigosPresentes    []string
	NumAmigosPresentes int
}

type ElementoListaNombresUsuario struct {
	Nombres []string
}

type ElementoListaUsuarios struct {
	Email           string
	NombreUsuario   string
	Biografia       string
	PartidasGanadas int
	PartidasTotales int
	Puntos          int
	ID_dado         int
	ID_ficha        int
}
