package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"strconv"
	"testing"
)

// Prueba de la funcionalidad de la fase de refuerzo de una partida desde 0
func TestFaseRefuerzoInicial(t *testing.T) {
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

	if len(partidaCache.Estado.Acciones) != (int(logica_juego.NUM_REGIONES) + 2) { // 42 regiones y acción de cambio de turno, cambio de fase
		t.Fatal("No se han asignado todas las regiones. Regiones asignadas:", len(partidaCache.Estado.Acciones))
	}

	numRegion := 0
	// Buscar región ocupada por usuario1
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		if partidaCache.Estado.EstadoMapa[i].Ocupante == "usuario1" {
			numRegion = int(i)
			break
		}
	}

	saltarTurnos(t, partidaCache, "usuario1")
	// Intentamos cambiar de fase con más de cuatro cartas, se espera error
	partidaCache = comprobarPartidaEnCurso(t, 1)

	numTropasRegionPrevio := partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(numRegion)].NumTropas
	numTropasPrevio := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas

	// Intentamos cambiar de fase con tropas no colocadas, se espera error
	t.Log("Intentamos cambiar de fase con tropas no colocadas, se espera error")
	err := saltarFase(cookie, t)
	partidaCache = comprobarPartidaEnCurso(t, 1)
	t.Log(partidaCache.Estado.EstadosJugadores["usuario1"].Tropas, partidaCache.Estado.Fase)
	if err == nil {
		t.Fatal("Se esperaba error al intentar saltar la fase")
	}
	t.Log("OK: No se ha podido saltar la fase, error:", err)

	// Reforzar con todas las tropas disponibles
	reforzarTerritorio(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas)

	t.Log("Intentamos cambiar de fase con 6 cartas, se espera error")
	partidaCache, err = cambioDeFaseConDemasiadasCartas(t, partidaCache, err, cookie)

	numTropasPost := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas
	numTropasRegionPost := partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(numRegion)].NumTropas

	if numTropasPrevio <= numTropasPost || numTropasPost != 0 {
		t.Fatal("Números de tropas incorrecto al agotarlas en una región. Prev:" + strconv.Itoa(numTropasPrevio) + "Post:" + strconv.Itoa(numTropasPost))
	} else if numTropasRegionPrevio >= numTropasRegionPost || numTropasRegionPost != (numTropasRegionPrevio+numTropasPrevio) {
		t.Fatal("Números de tropas en región incorrecto al agotarlas en una región. Prev:" + strconv.Itoa(numTropasRegionPrevio) + "Post:" + strconv.Itoa(numTropasRegionPost))
	} else {
		t.Log("Tropas asignadas correctamente. Tropas post:", numTropasPost)
	}

	// saltarTurnos(t, partidaCache, "usuario1")
	// Forzar fallo por no tener tropas
	reforzarTerritorioConFallo(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas+1)

	// Cambio de fase correcto
	t.Log("Intentamos cambiar de fase con todas las tropas colocadas")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error al saltar de fase:", err)
	}
	t.Log("Se ha saltado la fase correctamente")

	t.Log("Intentamos finalizar el ataque con territorios vacíos, se espera error")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.Egypt].Ocupante = ""
	globales.CachePartidas.AlmacenarPartida(partidaCache)
	err = saltarFase(cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al saltar el ataque con territorios vacíos")
	}
	t.Log("OK no se ha podido saltar el ataque con territorios vacíos, error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.Egypt].Ocupante = "Jugador1"
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Intentamos cambiar de fase con más de cuatro cartas, se espera error
	t.Log("Intentamos cambiar de ataque a fortificación con más de 4 cartas, se espera error")
	partidaCache, err = cambioDeFaseConDemasiadasCartas(t, partidaCache, err, cookie)

	// Cambio de fase correcto
	t.Log("Intentamos cambiar de fase, de ataque a fortificación")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error al saltar de fase:", err)
	}
	t.Log("Se ha saltado la fase correctamente")

	// Cambio de fase correcto
	t.Log("Intentamos cambiar de fase, de fortificación a refuerzo, cediendo el turno")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error al saltar de fase:", err)
	}
	t.Log("Se ha saltado la fase correctamente")

	//saltarTurnos(t, partidaCache, "usuario2")
	// Forzar fallo por estar fuera de turno
	reforzarTerritorioConFallo(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	t.Log("Acciones al final:", partidaCache.Estado.Acciones)
}
