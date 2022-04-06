package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

// Prueba de la funcionalidad de inicio de una partida desde 0
func TestInicioPartida(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuarios...")
	cookie := crearUsuario("usuario1", t)
	cookie2 := crearUsuario("usuario2", t)
	cookie3 := crearUsuario("usuario3", t)
	cookie4 := crearUsuario("usuario4", t)
	cookie5 := crearUsuario("usuario5", t)
	cookie6 := crearUsuario("usuario6", t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, true)

	partidas := obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas.")
	}

	partidaCache := comprobarPartidaNoEnCurso(t, 1)

	t.Log("Uniéndose a partida...")
	unirseAPartida(cookie2, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie3, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie4, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie5, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie6, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 0 {
		t.Fatal("Hay partidas, aunque ya ha tenido que empezar:", partidas)
	}
	partidaCache = comprobarPartidaEnCurso(t, 1)

	if len(partidaCache.Estado.Acciones) != (int(logica_juego.NUM_REGIONES) + 2) { // 42 regiones y acción de cambio de turno y cambio de fase
		t.Fatal("No se han asignado todas las regiones. Regiones asignadas:", len(partidaCache.Estado.Acciones))
	}

	for _, jugador := range partidaCache.Estado.Jugadores {
		t.Log("Estado de ", jugador, ":", partidaCache.Estado.EstadosJugadores[jugador])
	}

	comprobarAcciones(t, cookie)
	comprobarAcciones(t, cookie2)
	comprobarAcciones(t, cookie3)
	comprobarAcciones(t, cookie4)
	comprobarAcciones(t, cookie5)
	comprobarAcciones(t, cookie6)
}
