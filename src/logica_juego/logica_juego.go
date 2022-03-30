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
// que se debe aumentar en una unidad para
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
	Descartes  []Carta
	NumCambios int

	// Recibir carta
	HaConquistado   bool // True si ha conquistado algún territorio en el turno
	HaRecibidoCarta bool // True si ya ha robado carta, para evitar más de un robo

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
		Descartes:        []Carta{},
		NumCambios:       0,
		HaConquistado:    false,
		HaRecibidoCarta:  false,
	}

	return e
}

func (e *EstadoPartida) SiguienteJugadorSinAccion() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)
}

func (e *EstadoPartida) SiguienteJugador() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)

	e.Acciones = append(e.Acciones, AccionCambioTurno{IDAccion: 1, Jugador: e.Jugadores[e.TurnoJugador]})
}

// TODO: documentar
// TODO: Más elaborado (pseudo-random, con recorridos por el grafo, etc.)
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

// TODO documentar
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

	// TODO limitar refuerzos a la fase de refuerzo

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

// RecibirCarta da una carta a un jugador en caso de que haya conquistado un territorio durante su turno
func (e *EstadoPartida) RecibirCarta(jugador string) error {
	// Comprobamos que el jugador está en la partida y es su turno
	estado, existe := e.EstadosJugadores[jugador]
	if !existe {
		return errors.New("El jugador indicado," + jugador + ", no está en la partida")
	} else if !e.esTurnoJugador(jugador) {
		return errors.New("Se ha solicitado una acción fuera de turno, el jugador en este turno es " + e.ObtenerJugadorTurno())
	}

	if e.Fase != Fortificar {
		return errors.New("Solo se puede recibir una carta en la fase de fortificación")
	}

	if !e.HaConquistado {
		return errors.New("Se debe conquistar algún territorio para recibir una carta")
	}

	if e.HaRecibidoCarta {
		return errors.New("Solo se puede recibir una carta por turno")
	}

	e.HaRecibidoCarta = true
	carta, err := retirarPrimeraCarta(e.Cartas)
	if err != nil {
		// No quedan cartas en la baraja
		// Devolvemos los descartes a la baraja y barajamos
		copy(e.Cartas, e.Descartes)
		e.Descartes = nil
		barajarCartas(e.Cartas)
		carta, _ = retirarPrimeraCarta(e.Cartas)
	}
	estado.Cartas = append(estado.Cartas, carta)

	// Añadimos la acción
	e.Acciones = append(e.Acciones, NewAccionObtenerCarta(carta, jugador))
	return nil
}

// CambiarCartas permite al jugador cambiar un conjunto de 3 cartas por ejércitos
// Para ello, las cartas deberán ser del mismo tipo o cada una de un tipo diferente
// Además, se podrá realizar un cambio en caso de tener dos cartas del mismo tipo y un comodín
// Los cambios se realizarán durante la fase de refuerzo, o en fase de ataque, si el jugador tiene más
// de 4 cartas tras derrotar a un rival.
// Si alguno de los territorios de las cartas cambiadas están ocupados por el jugador, recibirá tropas extra
// El número de ejércitos recibidos dependerá del número total de canjes.
// Del primer al quinto cambio numTropas = 4 + (nº cambio - 1) * 2
// A partir del sexto cambio numTropas = 15 + (nº cambio - 6) * 5
func (e *EstadoPartida) CambiarCartas(jugador string, ID_carta1, ID_carta2, ID_carta3 int) error {
	// Comprobamos que el jugador está en la partida y es su turno
	estado, existe := e.EstadosJugadores[jugador]
	if !existe {
		return errors.New("El jugador indicado," + jugador + ", no está en la partida")
	} else if !e.esTurnoJugador(jugador) {
		return errors.New("Se ha solicitado una acción fuera de turno, el jugador en este turno es " + e.ObtenerJugadorTurno())
	}

	if e.Fase == Fortificar || (e.Fase == Ataque && len(estado.Cartas) < 5) {
		return errors.New("Solo se pueden cambiar cartas durante el refuerzo o el ataque," +
			" en caso de tener más de 5 tras derrotar a un rival")
	}

	if !existeCarta(ID_carta1, estado.Cartas) || !existeCarta(ID_carta2, estado.Cartas) ||
		!existeCarta(ID_carta3, estado.Cartas) {
		return errors.New("El jugador no dispone de todas las cartas para el cambio")
	}

	numeroCartasInicial := len(estado.Cartas)

	// Obtenemos las 3 cartas de la mano del jugador
	carta1, _ := retirarCartaPorID(ID_carta1, estado.Cartas)
	carta2, _ := retirarCartaPorID(ID_carta2, estado.Cartas)
	carta3, _ := retirarCartaPorID(ID_carta3, estado.Cartas)

	if !esCambioValido([]Carta{carta1, carta2, carta3}) {
		// Devolvemos las 3 cartas a la mano del jugador
		estado.Cartas = append(estado.Cartas, carta1, carta2, carta3)
		return errors.New("Las cartas introducidas no son válidas para realizar un cambio")
	}

	// Descartamos las 3 cartas
	e.Descartes = append(e.Descartes, carta1, carta2, carta3)

	// Calculamos el número de tropas a asignar
	numTropas := 0
	e.NumCambios++
	if e.NumCambios < 6 {
		numTropas += 4 + (e.NumCambios-1)*2
	} else {
		// Número de cambios >= 6
		numTropas += 15 + (e.NumCambios-6)*5
	}
	estado.Tropas += numTropas

	// TODO en caso de que haya varias regiones que coincidan, el jugador debería poder elegir a que región asignar los dos ejércitos extra
	hayBonificacion := false
	var regionBonificacion NumRegion

	regiones := obtenerRegionesCartas(estado.Cartas)
	for _, r := range regiones {
		if e.EstadoMapa[r].Ocupante == jugador {
			e.EstadoMapa[r].NumTropas += 2
			hayBonificacion = true
			regionBonificacion = r
			break
		}
	}

	// TODO cambiar accion de cambio de cartas -> se cambian conjuntos de uno en uno, numConjuntos no necesario

	e.Acciones = append(e.Acciones, NewAccionCambioCartas(1, numTropas, hayBonificacion, regionBonificacion, numeroCartasInicial >= 5))
	return nil
}

func (e *EstadoPartida) esTurnoJugador(jugador string) bool {
	return e.Jugadores[e.TurnoJugador] == jugador
}

func (e *EstadoPartida) ObtenerJugadorTurno() string {
	return e.Jugadores[e.TurnoJugador]
}
