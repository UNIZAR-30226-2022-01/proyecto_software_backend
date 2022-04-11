package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

func TestAtaqueUnitario(t *testing.T) {
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

	// Intento atacar con 5 o más cartas
	partida.Fase = logica_juego.Ataque
	partida.EstadosJugadores["Jugador1"].Cartas = []logica_juego.Carta{{IdCarta: 1}, {IdCarta: 1}, {IdCarta: 1},
		{IdCarta: 1}, {IdCarta: 1}, {IdCarta: 1}}
	t.Log("Intentamos atacar con más de 4 cartas, se espera error")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
	if err == nil {
		t.Fatal("Se esperaba error al atacar con más de 4 cartas")
	}
	t.Log("OK, se ha obtenido el error:", err)
	partida.EstadosJugadores["Jugador1"].Cartas = nil

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
	regionOrigen.NumTropas = 10
	t.Log("Intentamos realizar un ataque correcto")
	err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
	if err != nil {
		t.Fatal("Se ha obtenido el siguiente error al atacar:", err)
	}

	ultimaAccion := partida.Acciones[len(partida.Acciones)-1]
	ultimoAtaque, ok := ultimaAccion.(logica_juego.AccionAtaque)
	if !ok {
		t.Fatal("La última acción no es un ataque")
	}

	if ultimoAtaque.NumDadosAtaque != 3 {
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
		ultimoAtaque.NumDadosAtaque, "dados y ha perdido", ultimoAtaque.TropasPerdidasAtacante, "tropas")
	t.Log("El defensor", ultimoAtaque.JugadorDefensor, "ha perdido", ultimoAtaque.TropasPerdidasDefensor, "tropas")

	// Comprobamos el fin del ataque en caso de que el defensor se quede sin tropas
	partida.EstadoMapa[logica_juego.Venezuela].NumTropas = 10
	partida.EstadoMapa[logica_juego.Brazil].NumTropas = 1
	tropasDefensor := 1

	t.Log("Atacamos hasta que el defensor se quede sin tropas")
	for tropasDefensor > 0 {
		err = partida.Ataque(logica_juego.Venezuela, logica_juego.Brazil, 3, "Jugador1")
		if err != nil {
			t.Fatal("Se ha obtenido el siguiente error al atacar:", err)
		}

		ultimoAtaque, ok = partida.Acciones[len(partida.Acciones)-1].(logica_juego.AccionAtaque)
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
}
