// Package logica_juego define estructuras y constantes para los diferentes elementos del juego,
// así como la lista de acciones y regiones a ser interpretada por los clientes.
package logica_juego

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var onlyOnce sync.Once

// LanzarDados devuelve un número entre [0-5] (utilizado para indexar,
// que se debe aumentar en una unidad a la hora de mostrarse al usuario)
func LanzarDados() int {
	var dados = []int{0, 1, 2, 3, 4, 5}

	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	return dados[rand.Intn(len(dados))] // Devuelve una posición [0, 6)
}

type EstadoJugador struct {
	Cartas            []Carta
	UltimoIndiceLeido int
	Tropas            int
}

type EstadoRegion struct {
	Ocupante  string
	NumTropas int
}

type EstadoPartida struct {
	Acciones         []interface{} // Acciones realizadas durante la partida
	Jugadores        []string
	EstadosJugadores map[string]*EstadoJugador // Mapa de nombres de los jugadores en la partida y sus estados
	TurnoJugador     int                       // Índice de la lista que corresponde a qué jugador le toca

	Fase        Fase
	NumeroTurno int

	EstadoMapa map[NumRegion]*EstadoRegion

	// Baraja
	Cartas     []Carta
	NumCambios int

	// ...
}

func CrearEstadoPartida(jugadores []string) (e EstadoPartida) {
	e = EstadoPartida{
		Acciones:         make([]interface{}, 0),
		Jugadores:        crearSliceJugadores(jugadores),
		EstadosJugadores: crearMapaEstadosJugadores(jugadores),
		TurnoJugador:     LanzarDados(), // Primer jugador aleatorio
		Fase:             Inicio,
		NumeroTurno:      0,
		EstadoMapa:       crearEstadoMapa(),
		Cartas:           crearBaraja(),
		NumCambios:       0,
	}

	return e
}

// SiguienteJugadorSinAccion cambia el turno a otro jugador sin emitir ninguna acción.
// Usar únicamente al rellenar regiones.
func (e *EstadoPartida) SiguienteJugadorSinAccion() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)
}

// SiguienteJugador cambia el turno a otro jugador, emitiendo la acción correspondiente.
func (e *EstadoPartida) SiguienteJugador() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)

	// TODO: Encolar un AccionInicioTurno calculando el nº de tropas según las regiones y continentes ocupados
	//e.Acciones = append(e.Acciones, AccionCambioTurno{IDAccion: 1, Jugador: e.Jugadores[e.TurnoJugador]})
	e.Acciones = append(e.Acciones, struct{}{})
}

// RellenarRegiones rellena las regiones del estado de la partida equitativa y aleatoriamente entre los usuarios, consumiendo
// una tropa por cada región para controlarla. Aunque se emite una acción de asignación de territorio por cada uno de ellos
// asignado, no se emiten acciones de cambio de turno durante el proceso.
//
// Una vez terminado el proceso, se emite una acción de cambio de turno a un nuevo jugador.
// TODO: Más elaborado (pseudo-random teniendo en cuenta adyacencias, recorridos por el grafo, etc.)
func (e *EstadoPartida) RellenarRegiones() {
	regionesAsignadas := 0
	for i := Eastern_australia; i <= Alberta; i++ {
		if e.EstadosJugadores[e.Jugadores[e.TurnoJugador]].Tropas >= 1 {
			e.EstadoMapa[i].Ocupante = e.Jugadores[e.TurnoJugador]
			e.EstadoMapa[i].NumTropas = 1

			e.EstadosJugadores[e.Jugadores[e.TurnoJugador]].Tropas = e.EstadosJugadores[e.Jugadores[e.TurnoJugador]].Tropas - 1

			regionesAsignadas = regionesAsignadas + 1

			// Añadir una nueva acción
			e.Acciones = append(e.Acciones, NewAccionRecibirRegion(i, e.EstadosJugadores[e.Jugadores[e.TurnoJugador]].Tropas, NUM_REGIONES-regionesAsignadas, e.Jugadores[e.TurnoJugador]))

			e.SiguienteJugadorSinAccion()
		} else {
			// Repite la iteración para el siguiente jugador
			e.SiguienteJugadorSinAccion()
			i = i - 1
		}
	}

	// Se empieza con un jugador nuevo, sigue siendo pseudo-aleatorio frente al "offset" de las regiones asignadas
	e.SiguienteJugador()
}

// ReforzarTerritorio refuerza un territorio dado su id con numTropas para un jugador dado.
// Si la acción tiene éxito, se emite una acción de refuerzo, o nada y se devuelve un error ya formateado en caso contrario.
//
// Se hacen checks de:
//		Región incorrecta
//		Región ocupada por otro jugador
//		Jugador fuera de turno
//		Jugador con tropas insuficientes
func (e *EstadoPartida) ReforzarTerritorio(idTerritorio int, numTropas int, jugador string) error {
	region, existe := e.EstadoMapa[NumRegion(idTerritorio)]
	if !existe {
		return errors.New("La región indicada," + strconv.Itoa(idTerritorio) + ", es inválida")
	}

	if region.Ocupante != jugador {
		return errors.New("La región indicada," + strconv.Itoa(idTerritorio) + ", tiene por ocupante a otro jugador: " + region.Ocupante)
	}

	estado, existe := e.EstadosJugadores[jugador]
	if !existe {
		return errors.New("El jugador indicado," + jugador + ", no está en la partida")
	} else if !e.esTurnoJugador(jugador) {
		return errors.New("Se ha solicitado una acción fuera de turno, el jugador en este turno es " + e.ObtenerJugadorTurno())
	}

	if estado.Tropas-numTropas < 0 {
		return errors.New("No tienes tropas suficientes para reforzar un territorio, tropas restantes: " + strconv.Itoa(estado.Tropas))
	} else {
		estado.Tropas = estado.Tropas - numTropas
		region.NumTropas = region.NumTropas + numTropas

		// Añadir una nueva acción
		e.Acciones = append(e.Acciones, NewAccionReforzar(jugador, NumRegion(idTerritorio), numTropas))

		return nil
	}
}

func (e *EstadoPartida) esTurnoJugador(jugador string) bool {
	return e.Jugadores[e.TurnoJugador] == jugador
}

func (e *EstadoPartida) ObtenerJugadorTurno() string {
	return e.Jugadores[e.TurnoJugador]
}
