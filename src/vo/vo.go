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

	// TODO: representar estado del juego y chat de la partida
	Mensajes interface{}
	Estado   interface{}
}
