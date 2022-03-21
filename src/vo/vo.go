package vo

import "net/http"

// Usuario es un objeto de usuario equivalente al del modelo de base de datos.
type Usuario struct {
	IdUsuario       int
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
