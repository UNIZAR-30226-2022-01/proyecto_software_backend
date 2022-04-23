package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"reflect"
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

	estadoPrimeraLlamada := comprobarAcciones(t, cookie)
	comprobarAcciones(t, cookie2)
	comprobarAcciones(t, cookie3)
	comprobarAcciones(t, cookie4)
	comprobarAcciones(t, cookie5)
	comprobarAcciones(t, cookie6)

	// Comprobar la llamada de obtención de todas las acciones hasta el momento
	accionesCompl := preguntarEstadoCompleto(t, cookie)

	if !reflect.DeepEqual(estadoPrimeraLlamada.Acciones, accionesCompl.Acciones) {
		t.Log("Las acciones obtenidas en la primera llamada a obtener acciones y en la llamada de obtener todas no coinciden")
		t.Log("Acciones de llamada normal:", estadoPrimeraLlamada.Acciones)
		t.Fatal("Acciones de llamada de estado completo:", accionesCompl.Acciones)
	}

	estadoPrimeraLlamada = preguntarEstado(t, cookie)
	accionesCompl = preguntarEstadoCompleto(t, cookie)
	if len(estadoPrimeraLlamada.Acciones) != 0 && len(accionesCompl.Acciones) != len(partidaCache.Estado.Acciones) {
		t.Fatal("Se esperaba que la lista de acciones normal fuera vacía y la de acciones completa fuera igual que antes. Acciones normal:", estadoPrimeraLlamada.Acciones, ", acciones completa:", accionesCompl.Acciones)
	}

	t.Log("Lista de acciones completa:", accionesCompl.Acciones)

	// Comprobación de la llamada de estar en una partida, estando en una
	if !jugandoEnPartida(cookie, t) {
		t.Fatal("A \"usuario1\" no se le indica como jugando en ninguna partida al preguntarlo")
	}

	cookie7 := crearUsuario("usuario7", t)
	// Comprobación de la llamada de estar en una partida, no estando en ninguna
	if jugandoEnPartida(cookie7, t) {
		t.Fatal("A \"usuario7\" se le indica como jugando en una partida al preguntarlo")
	}
}
