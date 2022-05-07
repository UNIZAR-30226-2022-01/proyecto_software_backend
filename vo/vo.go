// Package vo define estructuras de datos a utilizar para la comunicación con la base
// de datos y los clientes
package vo

import (
	//"backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
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
	ID_avatar       int
}

type Partida struct {
	IdPartida          int
	EsPublica          bool
	PasswordHash       string
	EnCurso            bool
	MaxNumeroJugadores int

	// TODO: representar chat de la partida
	Mensajes []Mensaje
	Estado   logica_juego.EstadoPartida
}

type Mensaje struct {
	// TODO: representar tipo mensaje
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
func (p *Partida) IniciarPartida(jugadores []string) {
	p.EnCurso = true
	p.Estado = logica_juego.CrearEstadoPartida(jugadores)

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

type ElementoListaUsuariosSimilares struct {
	Nombre  string
	EsAmigo bool
}

type ElementoListaUsuarios struct {
	Email           string
	NombreUsuario   string
	Biografia       string
	PartidasGanadas int
	PartidasTotales int
	Puntos          int
	ID_dado         int
	ID_avatar       int
	EsAmigo         bool
}

type ElementoRankingUsuarios struct {
	NombreUsuario   string
	PartidasGanadas int
	PartidasTotales int
}

type EstadoLobby struct {
	EnCurso          bool
	EsPublico        bool
	Jugadores        int
	MaxJugadores     int
	NombresJugadores []string
}

type ItemTienda struct {
	Id          int
	Nombre      string
	Descripcion string
	Precio      int
	Tipo        string
}

type ResumenPartida struct {
	Jugadores        []string
	TurnoJugador     string
	Fase             logica_juego.Fase
	Terminada        bool
	EstadosJugadores map[string]ResumenEstadoJugador
	Mapa             map[logica_juego.NumRegion]logica_juego.EstadoRegion
}

type ResumenEstadoJugador struct {
	NumCartas int
	Cartas    []logica_juego.Carta
	Tropas    int
	Expulsado bool
	Eliminado bool
}
