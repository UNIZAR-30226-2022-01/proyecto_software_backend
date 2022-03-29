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

// IniciarPartida marca una partida como iniciada y crea un estado para ella con los jugadores indicados, iniciando la
// primera fase (asignación de territorios) tras ello.
func (p *Partida) IniciarPartida(jugadores []Usuario) {
	p.EnCurso = true
	p.Estado = CrearEstadoPartida(jugadores)

	p.Estado.RellenarRegiones()
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
