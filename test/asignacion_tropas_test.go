package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

func TestAsignacionTropas(t *testing.T) {
	// TODO modificar prueba para usar cambio de turno y evitar que AsignarTropasRefuerzo sea pública
	// Casos a probar
	// Territorios < 12
	// Territorios >= 12
	// Ocupa algún continente

	t.Log("Purgando DB...")
	purgarDB()

	estadoPartida := logica_juego.CrearEstadoPartida([]string{"Jugador1", "Jugador2", "Jugador3", "Jugador4", "Jugador5", "Jugador6"})
	estadoPartida.TurnoJugador = 5
	estadoPartida.Fase = logica_juego.Fortificar
	estadoJugador := estadoPartida.EstadosJugadores["Jugador1"]

	// Desocupamos el mapa
	for region, _ := range estadoPartida.EstadoMapa {
		estadoPartida.EstadoMapa[region].Ocupante = ""
	}

	estadoJugador.Tropas = 0

	// Comprobamos que se asignan 3 ejércitos en caso de tener menos de 12 territorios
	t.Log("El jugador no ocupa ningún territorio")
	estadoPartida.AsignarTropasRefuerzo("Jugador1")
	if estadoJugador.Tropas != 3 {
		t.Fatal("El jugador debería tener 3 tropas pero tiene", estadoJugador.Tropas)
	}
	t.Log("El jugador ha recibido", estadoJugador.Tropas, "tropas al principio del turno")

	// El jugador ocupa Asia
	// Deberá recibir 11 ejércitos
	t.Log("El jugador ocupa Asia (7), con 12 territorios (4)")
	for _, region := range logica_juego.Continentes["Asia"].Regiones {
		estadoPartida.EstadoMapa[region].Ocupante = "Jugador1"
	}

	estadoJugador.Tropas = 0
	estadoPartida.AsignarTropasRefuerzo("Jugador1")
	if estadoJugador.Tropas != 11 {
		t.Fatal("El jugador debería tener 11 tropas pero tiene", estadoJugador.Tropas)
	}
	t.Log("El jugador ha recibido", estadoJugador.Tropas, "tropas al principio del turno")

	// El jugador ocupa Asia y Europa
	// Un total de 19 regiones
	// Deberá recibir 5 ejércitos por Europa, 7 por Asia y 6 por territorios ocupados
	// Total ejércitos = 18
	t.Log("El jugador ocupa Asia (7) y Europa (5), con 19 territorios (6)")
	for _, region := range logica_juego.Continentes["Europa"].Regiones {
		estadoPartida.EstadoMapa[region].Ocupante = "Jugador1"
	}

	estadoJugador.Tropas = 0
	estadoPartida.AsignarTropasRefuerzo("Jugador1")
	if estadoJugador.Tropas != 18 {
		t.Fatal("El jugador debería tener 18 tropas pero tiene", estadoJugador.Tropas)
	}
	t.Log("El jugador ha recibido", estadoJugador.Tropas, "tropas al principio del turno")
}
