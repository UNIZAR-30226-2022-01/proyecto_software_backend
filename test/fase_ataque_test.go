package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

// TestAtaqueUnitario prueba todos los diferentes casos correctos y de error de la función de ataque
func TestAtaqueUnitario(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()
	var err error
	partida := logica_juego.CrearEstadoPartida([]string{"Jugador1", "Jugador2", "Jugador3", "Jugador4", "Jugador5", "Jugador6"})
	partida.RellenarRegiones()

	regionOrigen := partida.EstadoMapa[logica_juego.Venezuela]
	regionDestino := partida.EstadoMapa[logica_juego.Brazil]
	regionOrigen.Ocupante = "Jugador1"
	regionOrigen.NumTropas = 10
	regionDestino.Ocupante = "Jugador2"
	regionDestino.NumTropas = 3

	// Intento atacar fuera de turno
	partida.TurnoJugador = 5
	partida.Fase = logica_juego.Ataque
	t.Log("Intentamos atacar en el turno de otro jugador, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 2, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar fuera de turno")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intento atacar fuera de fase
	partida.TurnoJugador = 0
	partida.Fase = logica_juego.Refuerzo
	t.Log("Intentamos atacar en la fase de refuerzo, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 2, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar fuera de la fase correspondiente")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intento atacar con un territorio sin ocupar
	partida.HayTerritorioDesocupado = true
	t.Log("Intentamos atacar con algún territorio desocupado, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Alberta, 2, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar con un territorio desocupado")
	}
	t.Log("OK, se ha obtenido el error:", err)
	partida.HayTerritorioDesocupado = false

	// Intento atacar un territorio no adyacente
	partida.EstadoMapa[logica_juego.Alberta].Ocupante = "Jugador3"
	t.Log("Intentamos atacar a un territorio no adyacente, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Alberta, 2, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar un territorio no adyacente")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intento atacar un territorio controlado por mi mismo
	partida.EstadoMapa[logica_juego.Peru].Ocupante = "Jugador1"
	t.Log("Intentamos atacar a un territorio controlado por mi mismo, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Peru, 2, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar un territorio controlado por mi mismo")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intento atacar con un número incorrecto de dados
	t.Log("Intentamos atacar con menos de 1 dado, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 0, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar con menos de un dado")
	}
	t.Log("OK, se ha obtenido el error:", err)

	t.Log("Intentamos atacar con más de 3 dados, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 5, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar con más de 3 dados")
	}
	t.Log("OK, se ha obtenido el error:", err)

	t.Log("Intentamos atacar sin tener al menos un ejército más que dados, se espera error")
	regionOrigen.NumTropas = 3
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar menos ejércitos que dados utilizados")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Probamos ataques correctos
	partida.TurnoJugador = 0
	partida.Fase = logica_juego.Ataque
	regionOrigen.NumTropas = 10
	t.Log("Intentamos realizar un ataque correcto")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
	if err != nil {
		t.Fatal("Se ha obtenido el siguiente error al atacar:", err)
	}

	ultimaAccion := partida.Acciones[len(partida.Acciones)-1]
	ultimoAtaque, ok := ultimaAccion.(logica_juego.AccionAtaque)
	if !ok {
		t.Fatal("La última acción no es un ataque:", ultimaAccion)
	}

	if len(ultimoAtaque.DadosAtacante) != 3 {
		t.Fatal("El número de dados lanzados no es correcto")
	}
	if ultimoAtaque.TropasPerdidasDefensor+ultimoAtaque.TropasPerdidasAtacante != 2 {
		t.Fatal("No se han comparado dos dados")
	}
	if ultimoAtaque.JugadorAtacante != "Jugador1" || ultimoAtaque.JugadorDefensor != "Jugador2" {
		t.Fatal("Los jugadores del ataque no son los correspondientes")
	}

	t.Log("Se ha realizado correctamente el ataque desde", ultimoAtaque.Origen, "hasta", ultimoAtaque.Destino)
	t.Log("El jugador atacante", ultimoAtaque.JugadorAtacante, "ha utilizado",
		len(ultimoAtaque.DadosAtacante), "dados y ha perdido", ultimoAtaque.TropasPerdidasAtacante, "tropas")
	t.Log("El defensor", ultimoAtaque.JugadorDefensor, "ha perdido", ultimoAtaque.TropasPerdidasDefensor, "tropas")

	// Comprobamos el fin del ataque en caso de que el defensor se quede sin tropas
	partida.EstadoMapa[logica_juego.Venezuela].NumTropas = 100
	partida.EstadoMapa[logica_juego.Brazil].NumTropas = 1
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		partida.EstadoMapa[i].Ocupante = "Jugador3"
	}
	partida.EstadoMapa[logica_juego.Venezuela].Ocupante = "Jugador1"
	partida.EstadoMapa[logica_juego.Brazil].Ocupante = "Jugador2"
	partida.EstadosJugadores["Jugador2"].Cartas = []logica_juego.Carta{{IdCarta: 1}, {IdCarta: 2}}
	tropasDefensor := 1

	t.Log("Atacamos hasta que el defensor se quede sin tropas")
	for partida.EstadoMapa[logica_juego.Brazil].NumTropas > 0 {
		err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
		if err != nil {
			t.Fatal("Se ha obtenido el siguiente error al atacar:", err)
		}

		// La última es una eliminación
		ultimoAtaque, ok = partida.Acciones[len(partida.Acciones)-2].(logica_juego.AccionAtaque)
		if !ok {
			t.Fatal("La penúltima acción no es de ataque", ultimaAccion)
		}

		tropasDefensor -= ultimoAtaque.TropasPerdidasDefensor
	}

	// Comprobamos que se haya marcado correctamente que hay un territorio desocupado
	if !partida.HayTerritorioDesocupado || partida.EstadoMapa[logica_juego.Brazil].NumTropas > 0 {
		t.Fatal("El territorio defensor no ha sido conquistado")
	}
	t.Log("OK, el territorio defensor ha perdido todas sus tropas")

	// El defensor se queda sin territorios, comprobamos que el atacante recibe sus cartas
	if len(partida.EstadosJugadores["Jugador2"].Cartas) > 0 {
		t.Fatal("No se le han quitado las cartas al defensor")
	}
	if len(partida.EstadosJugadores["Jugador1"].Cartas) != 2 {
		t.Fatal("El atacante no ha recibido las cartas del defensor")
	}
	t.Log("OK, el atacante ha recibido las cartas del defensor")

	// Comprobamos que el jugador ha sido derrotado
	t.Log(partida.JugadoresActivos)
	if partida.JugadoresActivos[1] {
		t.Fatal("El jugador no ha sido eliminado de la partida correctamente")
	}
	t.Log("El defensor ha sido eliminado de la partida correctamente")

	// Comprobamos que al pasar de turno no se tiene en cuenta al jugador eliminado

	partida.SiguienteJugador()
	// Normalmente no podremos saltar de turno así, habrá que hacerlo pasando de fase
	// Ahora no se podría pasar de fase porque hay territorios sin ocupar

	if partida.TurnoJugador == 1 {
		t.Fatal("No se ha saltado al jugador eliminado")
	}
	t.Log("Se ha saltado correctamente al jugador eliminado")
}

// TestOcupacionUnitario prueba todos los diferentes casos correctos y de error de la función de ocupar
func TestOcupacionUnitario(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()
	var err error
	partida := logica_juego.CrearEstadoPartida([]string{"Jugador1", "Jugador2", "Jugador3", "Jugador4", "Jugador5", "Jugador6"})
	partida.RellenarRegiones()

	regionOrigen := partida.EstadoMapa[logica_juego.Venezuela]
	regionDestino := partida.EstadoMapa[logica_juego.Brazil]
	regionOrigen.Ocupante = "Jugador1"
	regionOrigen.NumTropas = 100
	regionDestino.Ocupante = "Jugador2"
	regionDestino.NumTropas = 1

	// Atacamos hasta que el defensor se quede sin tropas
	partida.TurnoJugador = 0
	partida.Fase = logica_juego.Ataque
	t.Log("Atacamos hasta que el defensor se quede sin tropas")
	tropasDefensor := 1
	for tropasDefensor > 0 {
		err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
		if err != nil {
			t.Fatal("Se ha obtenido el siguiente error al atacar:", err)
		}

		ultimoAtaque, ok := partida.Acciones[len(partida.Acciones)-1].(logica_juego.AccionAtaque)
		if !ok {
			t.Fatal("La última acción no es de ataque")
		}

		tropasDefensor -= ultimoAtaque.TropasPerdidasDefensor
	}

	// Comprobamos que se haya marcado correctamente que hay un territorio desocupado
	if !partida.HayTerritorioDesocupado || partida.EstadoMapa[logica_juego.Brazil].NumTropas > 0 {
		t.Fatal("El territorio defensor no ha sido conquistado")
	}
	t.Log("OK, el territorio defensor ha perdido todas sus tropas")

	// Comprobamos que el último ataque sea correcto
	ultimoAtaque, ok := partida.Acciones[len(partida.Acciones)-1].(logica_juego.AccionAtaque)
	if !ok {
		t.Fatal("La última acción no es de ataque")
	}

	if ultimoAtaque.JugadorDefensor != "Jugador2" {
		t.Fatal("El jugador defensor del último ataque no es correcto")
	}
	if partida.RegionUltimoAtaque != logica_juego.Venezuela {
		t.Fatal("La región del último ataque no es correcta")
	}
	t.Log("OK, el último ataque ha sido correcto")

	// Intentamos ocupar fuera de turno
	partida.TurnoJugador = 5
	partida.Fase = logica_juego.Ataque
	t.Log("Intentamos ocupar fuera de turno, se espera error")
	err = partida.Ocupar(logica_juego.Brazil, 20, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar fuera de turno")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intentamos ocupar fuera de fase
	partida.TurnoJugador = 0
	partida.Fase = logica_juego.Refuerzo
	t.Log("Intentamos ocupar en una fase incorrecta, se espera error")
	err = partida.Ocupar(logica_juego.Brazil, 20, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar fuera de la fase de ataque")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intentamos ocupar con más de 4 cartas
	partida.Fase = logica_juego.Ataque
	partida.EstadosJugadores["Jugador1"].Cartas = []logica_juego.Carta{{IdCarta: 1}, {IdCarta: 1}, {IdCarta: 1},
		{IdCarta: 1}, {IdCarta: 1}, {IdCarta: 1}}
	t.Log("Intentamos ocupar con más de 5 cartas, se espera error")
	err = partida.Ocupar(logica_juego.Brazil, 20, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar con más de 5 cartas")
	}
	t.Log("OK, se ha obtenido el error:", err)

	partida.EstadosJugadores["Jugador1"].Cartas = nil

	// Intentamos ocupar con el estado indicando que no hay territorios desocupados
	partida.HayTerritorioDesocupado = false
	t.Log("Intentamos ocupar sin territorios desocupados, se espera error")
	err = partida.Ocupar(logica_juego.Brazil, 20, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar sin territorios desocupados")
	}
	t.Log("OK, se ha obtenido el error:", err)
	partida.HayTerritorioDesocupado = true

	// Intentamos ocupar un territorio con tropas
	partida.EstadoMapa[logica_juego.Peru].NumTropas = 3
	partida.EstadoMapa[logica_juego.Peru].Ocupante = "Jugador3"
	t.Log("Intentamos ocupar un territorio con tropas, se espera error")
	err = partida.Ocupar(logica_juego.Peru, 20, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar un territorio con tropas")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intentamos ocupar un territorio no adyacente
	partida.EstadoMapa[logica_juego.Alberta].NumTropas = 0
	partida.EstadoMapa[logica_juego.Alberta].Ocupante = "Jugador3"
	t.Log("Intentamos ocupar un territorio no adyacente, se espera error")
	err = partida.Ocupar(logica_juego.Alberta, 20, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar un territorio no adyacente")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intentamos ocupar con un número de tropas menor a numDados - numEjercitosPerdidos
	numDados := partida.DadosUltimoAtaque
	numEjercitosPerdidos := partida.TropasPerdidasUltimoAtaque
	t.Log("Intentamos ocupar un territorio utilizando menos ejércitos que el número de dados de la última tirada,\n" +
		" menos el número de ejércitos perdidos en el último ataque, se espera error")
	err = partida.Ocupar(logica_juego.Brazil, numDados-numEjercitosPerdidos-1, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar con tropas < dadosUltimoAtaque - tropasPerdidasUltimoAtaques")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intentamos ocupar dejando el territorio de origen sin tropas
	t.Log("Intentamos ocupar un territorio dejando el territorio original sin tropas, se espera error")
	err = partida.Ocupar(logica_juego.Brazil, 120, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al ocupar dejando el territorio original sin tropas")
	}
	t.Log("OK, se ha obtenido el error:", err)

	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		partida.EstadoMapa[i].Ocupante = "Jugador1"
	}
	partida.EstadoMapa[logica_juego.Brazil].Ocupante = "Jugador2"

	// Realizamos una ocupación correcta
	t.Log("Intentamos realizar una ocupación correcta")
	partida.EstadoMapa[logica_juego.Venezuela].NumTropas = 100
	err = partida.Ocupar(logica_juego.Brazil, 40, "Jugador1")
	if err != nil {
		t.Fatal("Error al ocupar:", err)
	}
	ultimaOcupacion, ok := partida.Acciones[len(partida.Acciones)-1].(logica_juego.AccionOcupar)

	// Comprobamos la corrección de la ocupación
	if !ok {
		t.Fatal("La última acción no es una ocupación")
	}
	t.Log("Tropas origen:", ultimaOcupacion.TropasOrigen, ", Tropas destino:", ultimaOcupacion.TropasDestino)
	if ultimaOcupacion.JugadorOcupado != "Jugador2" || ultimaOcupacion.JugadorOcupante != "Jugador1" {
		t.Fatal("Los jugadores de la última ocupación no son correctos")
	}
	if ultimaOcupacion.Destino != logica_juego.Brazil || ultimaOcupacion.Origen != logica_juego.Venezuela {
		t.Fatal("Los territorios de la última ocupación no son correctos")
	}
	if partida.HayTerritorioDesocupado {
		t.Fatal("No se ha actualizado correctamente la ocupación en el estado de la partida")
	}
	if ultimaOcupacion.TropasOrigen != 60 || ultimaOcupacion.TropasDestino != 40 {
		t.Fatal("El número de tropas tras la ocupación no es correcto")
	}

	t.Log("OK, se ha realizado la ocupación correctamente")

}

// TestIntegracionAtaqueOcupar simula una fase de ataque de una partida, probando diferentes casos correctos y
// de error de las acciones de ataque y ocupación
func TestIntegracionAtaqueOcupar(t *testing.T) {
	var err error
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

	// Modificar estado para pasar turno y a fase de ataque
	partidaCache := comprobarPartidaEnCurso(t, 1)
	saltarTurnos(t, partidaCache, "usuario5")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.Fase = logica_juego.Ataque
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Configuramos el mapa para las pruebas
	partidaCache.Estado.EstadoMapa[logica_juego.Venezuela].Ocupante = "usuario1"
	partidaCache.Estado.EstadoMapa[logica_juego.Venezuela].NumTropas = 100
	partidaCache.Estado.EstadoMapa[logica_juego.Peru].Ocupante = "usuario1"
	partidaCache.Estado.EstadoMapa[logica_juego.Brazil].Ocupante = "usuario2"
	partidaCache.Estado.EstadoMapa[logica_juego.Brazil].NumTropas = 1
	partidaCache.Estado.EstadoMapa[logica_juego.Alberta].Ocupante = "usuario2"
	partidaCache.Estado.EstadoMapa[logica_juego.Alberta].Ocupante = "usuario3"
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Intento atacar fuera de turno, se espera error
	t.Log("Intento atacar fuera de turno, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Brazil, 1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar fuera de turno")
	}
	t.Log("OK, se ha recibido el error:", err)

	// Saltar a turno de "usuario1" y fase de fortificación
	partidaCache = comprobarPartidaEnCurso(t, 1)
	saltarTurnos(t, partidaCache, "usuario1")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.Fase = logica_juego.Fortificar
	globales.CachePartidas.AlmacenarPartida(partidaCache)
	partidaCache = comprobarPartidaEnCurso(t, 1)

	// Intento atacar fuera de fase, se espera error
	t.Log("Intento atacar fuera de fase, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Brazil, 1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar fuera de fase")
	}
	t.Log("OK, se ha recibido el error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.Fase = logica_juego.Ataque
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Intento atacar un territorio controlado por mi mismo, se espera error
	t.Log("Intento atacar un territorio propio, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Peru, 1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar un territorio propio")
	}
	t.Log("OK, se ha recibido el error:", err)

	// Intento atacar un territorio no aydacente, se espera error
	t.Log("Intento atacar un territorio no adyacente, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Alberta, 1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar un territorio no adyacente")
	}
	t.Log("OK, se ha recibido el error:", err)

	// Intento atacar desde un territorio que no me pertenece, se espera error
	t.Log("Intento atacar desde un territorio que no me pertenece, se espera error")
	err = atacar(logica_juego.Argentina, logica_juego.Brazil, 1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar un territorio que no me pertenece")
	}
	t.Log("OK, se ha recibido el error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadosJugadores["usuario1"].Cartas = nil
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Intento atacar con algún territorio sin ocupar, se espera error
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.HayTerritorioDesocupado = true
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	t.Log("Intento atacar con algún territorio sin ocupar, se espera error")
	err = atacar(logica_juego.Argentina, logica_juego.Brazil, 1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar con algún territorio sin ocupar")
	}
	t.Log("OK, se ha recibido el error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.HayTerritorioDesocupado = false
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Intento atacar con un número de dados incorrecto, se espera error
	t.Log("Intento atacar con menos de 1 dado, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Brazil, 0, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar con menos de 1 dado")
	}
	t.Log("OK, se ha recibido el error:", err)

	t.Log("Intento atacar con más de 3 dados, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Brazil, 4, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar con más de 3 dados")
	}
	t.Log("OK, se ha recibido el error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.Venezuela].NumTropas = 2
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	t.Log("Intento atacar sin tener más ejércitos que dados, se espera error")
	err = atacar(logica_juego.Venezuela, logica_juego.Brazil, 2, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al atacar sin tener más ejércitos que dados")
	}
	t.Log("OK, se ha recibido el error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.Venezuela].NumTropas = 100
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Atacamos hasta conquistar el territorio
	t.Log("Intento conquistar un territorio")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	for !partidaCache.Estado.HayTerritorioDesocupado {
		t.Log("Ataco al territorio")
		err = atacar(logica_juego.Venezuela, logica_juego.Brazil, 3, cookie, t)
		if err != nil {
			t.Fatal("Error al atacar:", err)
		}

		partidaCache = comprobarPartidaEnCurso(t, 1)
	}
	t.Log("He conquistado")

	numDadosUltimoAtaque := partidaCache.Estado.DadosUltimoAtaque
	numTropasPerdidas := partidaCache.Estado.TropasPerdidasUltimoAtaque
	minimoTropas := numDadosUltimoAtaque - numTropasPerdidas

	// Intento ocupar el territorio con menos tropas de las necesarias, se espera error
	t.Log("Intento ocupar con menos tropas de las necesarias, se espera error")
	err = ocupar(logica_juego.Brazil, minimoTropas-1, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al ocupar sin con menos tropas de las necesarias")
	}
	t.Log("OK, se ha recibido el error:", err)

	// Intento ocupar el territorio dejando el origen sin tropas, se espera error
	t.Log("Intento ocupar dejando al origen sin tropas, se espera error")
	err = ocupar(logica_juego.Brazil, 150, cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al ocupar dejando al origen sin tropas")
	}
	t.Log("OK, se ha recibido el error:", err)

	// Intento cambiar de fase sin ocupar
	t.Log("Intento cambiar de fase sin ocupar, se espera error")
	err = saltarFase(cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al saltar de fase dejando territorios sin ocupar")
	}
	t.Log("OK, se ha recibido el error:", err)

	// Ocupamos el territorio
	t.Log("Intento ocupar el territorio")
	err = ocupar(logica_juego.Brazil, 50, cookie, t)
	if err != nil {
		t.Fatal("Error al ocupar:", err)
	}
	t.Log("OK, se ha ocupado el territorio correctamente")

	// Cambio de fase tras la ocupación
	// Intento cambiar de fase tras ocupar
	t.Log("Intento cambiar de fase tras ocupar")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error obtenido al cambiar de fase:", err)
	}
	partidaCache = comprobarPartidaEnCurso(t, 1)
	if partidaCache.Estado.Fase != logica_juego.Fortificar {
		t.Fatal("No se ha cambiado a la fase de fortificar correctamente")
	}
	t.Log("OK, se ha cambiado de fase correctamente")

	// Cambio de turno y espero a recibir una carta
	t.Log("Finalizo el turno")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error obtenido al cambiar de fase:", err)
	}

	// Comprobamos que la antepenúltima acción sea la carta
	// Acciones: ..., recibir carta, cambio de fase y cambio de turno
	t.Log("Comprobamos si el jugador ha recibido una carta por conquistar")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	acciones := partidaCache.Estado.Acciones
	accionRecibirCarta, ok := acciones[len(acciones)-3].(logica_juego.AccionObtenerCarta)
	if !ok {
		t.Fatal("El jugador no ha recibido una carta tras conquistar el territorio")
	}
	if accionRecibirCarta.Jugador != "usuario1" {
		t.Fatal("La carta no se ha otorgado al jugador correspondiente")
	}
	t.Log("OK, se ha introducio la siguiente accion de recibir carta", accionRecibirCarta)
}
