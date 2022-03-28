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
	Jugadores          []Usuario // TODO: Eliminar

	// TODO: representar estado del juego y chat de la partida
	Mensajes []Mensaje
	Estado   EstadoPartida
}

type Mensaje struct {
	// TODO
}

func CrearPartida(esPublica bool, passwordHash string, maxNumeroJugadores int) *Partida {
	partida := Partida{
		IdPartida:          0,
		EsPublica:          esPublica,
		PasswordHash:       passwordHash,
		EnCurso:            false,
		MaxNumeroJugadores: maxNumeroJugadores,
		//Estado: logica_juego.CrearEstadoPartida()
	}

	return &partida
}

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

// ContarAmigos devuelve cuántos amigos de un usuario están en una partida dada
func ContarAmigos(amigos []Usuario, partida Partida) (num int) {
	for _, amigo := range amigos {
		// Como máximo hay 6 jugadores en la partida, así que
		// la complejidad la dicta el número de amigos del usuario
		for _, jugador := range partida.Jugadores {
			if amigo.NombreUsuario == jugador.NombreUsuario {
				num = num + 1
			}
		}
	}

	return num
}
