package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"log"
	"testing"
)

func TestBaraja(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	estadoPartida := logica_juego.CrearEstadoPartida([]string{"Jugador1", "Jugador2", "Jugador3"})
	estadoPartida.TurnoJugador = 0
	estadoPartida.Fase = logica_juego.Fortificar
	robarBarajaCompleta(&estadoPartida, t)
	if len(estadoPartida.Cartas) > 0 {
		t.Fatal("La baraja debería estar vacía, pero tiene", len(estadoPartida.Cartas), "cartas")
	}
	estadoJugador := estadoPartida.EstadosJugadores["Jugador1"]
	t.Log("Baraja completa:", estadoJugador.Cartas)
	var err error
	var carta logica_juego.Carta
	var baraja []logica_juego.Carta

	// Prueba función retirar carta por ID
	// Retiramos todas las cartas de la mano del jugador
	for i := 0; i < 44; i++ {
		carta, estadoJugador.Cartas, err = logica_juego.RetirarCartaPorID(i, estadoJugador.Cartas)
		if err != nil {
			t.Fatal("Error al retirar cartas del jugador:", err)
		}
		if carta.IdCarta != i {
			t.Fatal("La carta recuperada no es correcta, id:", carta.IdCarta, "se esperaba: ", i)
		}
		baraja = append(baraja, carta)
	}
	t.Log("Cartas en la mano del jugador:", estadoJugador.Cartas)
	if len(estadoJugador.Cartas) > 0 {
		t.Fatal("El jugador no debería tener cartas")
	}

	// Compiamos las cartas al montón de descartes
	estadoPartida.Descartes = baraja

	// Tomamos una carta de la baraja
	// Al estar vacía, toma los descartes y los rebaraja
	estadoPartida.HaConquistado = true
	estadoPartida.HaRecibidoCarta = false

	err = estadoPartida.RecibirCarta("Jugador1")
	if err != nil {
		t.Fatal("Error al recibir carta:", err)
	}
	estadoPartida.HaConquistado = true
	estadoPartida.HaRecibidoCarta = false

	// La baraja debería tener 43 cartas
	if len(estadoPartida.Cartas) != 43 {
		t.Fatal("La baraja debería tener 43 cartas, pero tiene:", len(estadoPartida.Cartas))
	}

	// Volvemos a tomar toda la baraja
	robarBarajaCompleta(&estadoPartida, t)
	if len(estadoPartida.Cartas) > 0 {
		t.Fatal("La baraja debería estar vacía, pero tiene:", len(estadoPartida.Cartas), "cartas")
	}

	// Prueba de canjes
	// Cambiamos 3 cartas de infantería
	estadoPartida.Fase = logica_juego.Refuerzo
	t.Log("Cambiando 3 cartas de infanteria")
	cambiarCartas(t, estadoJugador, &estadoPartida, 0, 1, 2, 1)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos 3 cartas de caballeria
	t.Log("Cambiando 3 cartas de caballeria")
	cambiarCartas(t, estadoJugador, &estadoPartida, 18, 19, 20, 2)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos 3 cartas de artilleria
	t.Log("Cambiando 3 cartas de artilleria")
	cambiarCartas(t, estadoJugador, &estadoPartida, 36, 37, 38, 3)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos una de cada
	t.Log("Cambiando 3 cartas, una de cada tipo")
	cambiarCartas(t, estadoJugador, &estadoPartida, 3, 21, 39, 4)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiar con dos comodines
	t.Log("Intentamos cambiar cartas con dos comodines, se espera error")
	err = estadoPartida.CambiarCartas("Jugador1", 4, 42, 43)
	if err == nil {
		t.Fatal("Se esperaba obtener error al cambiar con 2 comodines")
	}

	t.Log("OK: No se ha podido cambiar cartas con dos comodines, error:", err)

	// Cambiamos 2 cartas de un tipo + un comodín
	t.Log("Cambiando 3 cartas, 2 de infantería y un comodín")
	cambiarCartas(t, estadoJugador, &estadoPartida, 4, 5, 42, 5)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Prueba de errores en el canje
	// Cambiar cartas que no tenemos
	t.Log("Intentamos cambiar cartas que no tenemos, se espera error")
	err = estadoPartida.CambiarCartas("Jugador1", 0, 1, 2)
	if err == nil {
		t.Fatal("Se esperaba obtener error al cambiar con cartas que el jugador no tiene")
	}

	t.Log("OK: No se ha podido cambiar con cartas que el jugador no tiene, error:", err)

	// Cambiar con cartas de distinto tipo
	t.Log("Intentamos cambiar con cartas de distinto tipo, se espera error")
	err = estadoPartida.CambiarCartas("Jugador1", 6, 7, 23)
	if err == nil {
		t.Fatal("Se esperaba obtener error al cambiar con cartas de distintos tipos")
	}

	t.Log("OK: No se ha podido cambiar con cartas de diferentes tipos, error:", err)

	// Probamos bonificación por territorio
	// TODO probar caso en el que más de una carta recibe bonificación por territorio
	carta, estadoJugador.Cartas, err = logica_juego.RetirarCartaPorID(10, estadoJugador.Cartas)
	if err != nil {
		t.Fatal("Error al tomar la carta de la mano del jugador")
	}
	estadoJugador.Cartas = append(estadoJugador.Cartas, carta)
	region := carta.Region
	estadoRegion := estadoPartida.EstadoMapa[region]
	estadoRegion.Ocupante = "Jugador1"
	tropasIniciales := estadoRegion.NumTropas
	t.Log("Cambiando 3 cartas, una de ellas con bonificación por territorio")
	cambiarCartas(t, estadoJugador, &estadoPartida, 10, 11, 12, 6)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)
	if estadoRegion.NumTropas-tropasIniciales != 2 {
		t.Fatal("No se han recibido tropas adicionales")
	}
	t.Log("Inicialmente, la región tenía", tropasIniciales, "tropas, tras el canje hay un total de ", estadoRegion.NumTropas-tropasIniciales)

	// Comprobar que se añade correctamente la acción de cambio de cartas
	accion := estadoPartida.Acciones[len(estadoPartida.Acciones)-1]
	accionCambio, ok := accion.(logica_juego.AccionCambioCartas)
	if !ok {
		t.Fatal("La última acción no es un cambio de cartas")
	}

	if !accionCambio.ObligadoAHacerCambios {
		t.Fatal("El jugador debería estar obligado a cambiar, al tener más de 4 cartas")
	}

	if !accionCambio.BonificacionObtenida {
		t.Fatal("No se ha obtenido bonificación por territorio al cambiar las cartas")
	}

	t.Log("El último cambio de cartas fue:", accionCambio)

	// Probar cambios de cartas de número > 6
	t.Log("Cambiando 3 cartas, 7º cambio")
	cambiarCartas(t, estadoJugador, &estadoPartida, 30, 31, 32, 7)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	t.Log("Cambiando 3 cartas, 8º cambio")
	cambiarCartas(t, estadoJugador, &estadoPartida, 33, 34, 35, 8)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Consultar cartas de un jugador

	t.Log("Creando usuarios...")
	cookie := crearUsuario("usuario1", t)
	cookie2 := crearUsuario("usuario2", t)
	cookie3 := crearUsuario("usuario3", t)
	cookie4 := crearUsuario("usuario4", t)
	cookie5 := crearUsuario("usuario5", t)
	cookie6 := crearUsuario("usuario6", t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, true)

	t.Log("Uniéndose a partida...")
	unirseAPartida(cookie2, t, 1)
	unirseAPartida(cookie3, t, 1)
	unirseAPartida(cookie4, t, 1)
	unirseAPartida(cookie5, t, 1)
	unirseAPartida(cookie6, t, 1)
	partidaCache := comprobarPartidaEnCurso(t, 1)
	cartasJugador := partidaCache.Estado.Cartas[0:3]
	partidaCache.Estado.EstadosJugadores["usuario1"].Cartas = cartasJugador
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	cartasObtenidas := consultarCartas(cookie, t)
	log.Println("Se deberían recibir las siguientes cartas:", cartasJugador)
	log.Println("Se han recibido estas cartas:", cartasObtenidas)
	if len(cartasObtenidas) != len(cartasJugador) {
		t.Fatal("No se ha recibido el mismo número de cartas")
	}

	for i, carta := range cartasJugador {
		if carta != cartasJugador[i] {
			t.Fatal("No se han recibido las mismas cartas")
		}
	}
}
