package vo

import "net/http"

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
	NumeroJugadores    int
	MaxNumeroJugadores int
	Jugadores          []Usuario

	// TODO: representar estado del juego y chat de la partida
	Mensajes interface{}
	Estado   interface{}
}

// ContarAmigos devuelve cuántos amigos de un usuario están en una partida dada
func ContarAmigos(amigos []Usuario, partida Partida) (num int) {
	for _, amigo := range amigos {
		// Como máximo hay 6 jugadores en la partida, así que
		// la complejidadlo dicta el número de amigos del usuario
		for _, jugador := range partida.Jugadores {
			if amigo.NombreUsuario == jugador.NombreUsuario {
				num = num + 1
			}
		}
	}

	return num
}
