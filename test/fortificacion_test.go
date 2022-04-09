package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

// Prueba de fortificación, forzando fortificaciones válidas e inválidas a
// partir de un estado de regiones específico
func TestFortificacion(t *testing.T) {
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
	unirseAPartida(cookie2, t, 1)
	unirseAPartida(cookie3, t, 1)
	unirseAPartida(cookie4, t, 1)
	unirseAPartida(cookie5, t, 1)
	unirseAPartida(cookie6, t, 1)

	partidaCache := comprobarPartidaEnCurso(t, 1)

	// Saltar a turno de "usuario1" y fase de fortificación
	saltarTurnos(t, partidaCache, "usuario1")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	pasarAFaseFortificar(partidaCache)

	// "Conquistar" regiones con algunas ocupadas en el camino, con el resto perteneciente a
	// otros usuarios, y conquistar otra desconectada
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		partidaCache = conquistar(t, partidaCache, int(i), "usuario2")
	}
	// Alaska -> Alberta -> Ontario -> Eastern US
	partidaCache = conquistar(t, partidaCache, int(logica_juego.Alaska), "usuario1")
	partidaCache = conquistar(t, partidaCache, int(logica_juego.Alberta), "usuario1")
	partidaCache = conquistar(t, partidaCache, int(logica_juego.Ontario), "usuario1")
	partidaCache = conquistar(t, partidaCache, int(logica_juego.Eastern_united_states), "usuario1")

	// Venezuela(desconectada del resto)
	partidaCache = conquistar(t, partidaCache, int(logica_juego.Venezuela), "usuario1")

	// Comprobar las regiones de su subgrafo
	var regiones []logica_juego.NumRegion
	grafo := partidaCache.Estado.ObtenerSubgrafoRegiones("usuario1")
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		if grafo.Node(int64(i)) != nil {
			regiones = append(regiones, i)
		}
	}

	if len(regiones) < 2 {
		t.Fatal("Se esperaban al menos dos (más) regiones controladas tras iniciar partida, obtenidas:", regiones)
	}

	t.Log("Territorios controlados:", regiones)

	// Dar 20 tropas al "usuario1"
	partidaCache = darTropas(t, partidaCache, 20, "usuario1")
	// Marca todos los territorios con 20 tropas
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		partidaCache.Estado.EstadoMapa[i].NumTropas = 20
	}
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Refuerzo válido
	t.Log("Reforzando territorio", logica_juego.Eastern_united_states, "desde", logica_juego.Alaska)
	fortificarTerritorio(t, cookie, 15, int(logica_juego.Alaska), int(logica_juego.Eastern_united_states))

	// Refuerzo inválido por regiones desconectadas
	t.Log("Reforzando territorio", logica_juego.Venezuela, "desde", logica_juego.Alaska)
	fortificarTerritorioConError(t, cookie, 2, int(logica_juego.Alaska), int(logica_juego.Venezuela))

	// Refuerzo inválido por números inválidos de tropas
	t.Log("Reforzando territorio", logica_juego.Eastern_united_states, "desde", logica_juego.Alaska)
	fortificarTerritorioConError(t, cookie, 20, int(logica_juego.Alaska), int(logica_juego.Eastern_united_states))

	// Refuerzo inválido por ser territorios controlados por otro
	t.Log("Reforzando territorio", logica_juego.Egypt, "desde", logica_juego.Eastern_united_states)
	fortificarTerritorioConError(t, cookie, 1, int(logica_juego.Eastern_united_states), int(logica_juego.Egypt))

	// Refuerzo inválido por ser 0 tropas
	t.Log("Reforzando territorio", logica_juego.Eastern_united_states, "desde", logica_juego.Alaska)
	fortificarTerritorioConError(t, cookie, 0, int(logica_juego.Alaska), int(logica_juego.Eastern_united_states))

	// Terminar la fase
	err := saltarFase(cookie, t)
	if err != nil {
		t.Fatal(err)
	}

	partidaCache = comprobarPartidaEnCurso(t, 1)
	if partidaCache.Estado.Jugadores[partidaCache.Estado.TurnoJugador] == "usuario1" {
		t.Fatal("No se ha cambiado de turno al pasar de fase desde fortificación")
	} else {
		t.Log("Fase terminada, turno actual:", partidaCache.Estado.Jugadores[partidaCache.Estado.TurnoJugador], ", fase:", partidaCache.Estado.Fase)
	}
}
