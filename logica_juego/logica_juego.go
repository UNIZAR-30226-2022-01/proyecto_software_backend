// Package logica_juego define estructuras y constantes para los diferentes elementos del juego,
// así como la lista de acciones y regiones a ser interpretada por los clientes.
package logica_juego

import (
	"errors"
	"log"
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
	Descartes  []Carta
	NumCambios int

	// Recibir carta
	HaConquistado   bool // True si ha conquistado algún territorio en el turno
	HaRecibidoCarta bool // True si ya ha robado carta, para evitar más de un robo
	HaFortificado   bool // True si ya ha fortificado en el turno

	// ...
}

func CrearEstadoPartida(jugadores []string) (e EstadoPartida) {
	e = EstadoPartida{
		Acciones:         make([]interface{}, 0),
		Jugadores:        crearSliceJugadores(jugadores),
		EstadosJugadores: crearMapaEstadosJugadores(jugadores),
		TurnoJugador:     (LanzarDados()) % len(jugadores), // Primer jugador aleatorio
		Fase:             Inicio,
		NumeroTurno:      0,
		EstadoMapa:       crearEstadoMapa(),
		Cartas:           crearBaraja(),
		Descartes:        []Carta{},
		NumCambios:       0,
		HaConquistado:    false,
		HaRecibidoCarta:  false,
		HaFortificado:    false,
	}

	return e
}

// SiguienteJugadorSinAccion cambia el turno a otro jugador sin emitir ninguna acción.
// Usar únicamente al rellenar regiones.
func (e *EstadoPartida) SiguienteJugadorSinAccion() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)
}

// SiguienteJugador cambia el turno a otro jugador, emitiendo la acción correspondiente.
// TODO no debería ser pública, lo es por necesidad dentro del test
func (e *EstadoPartida) SiguienteJugador() {
	// TODO cuando se puedan eliminar jugadores, habrá que tenerlo en cuenta a la hora de cambiar de turno
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)

	// En el nuevo turno no se habrá recibido carta, conquistado ni fortificado
	e.HaRecibidoCarta = false
	e.HaConquistado = false
	e.HaFortificado = false

	// Pasamos a la fase de refuerzo
	e.Fase = Refuerzo

	if e.Fase != Inicio {
		e.Acciones = append(e.Acciones, NewAccionCambioFase(Refuerzo, e.Jugadores[e.TurnoJugador]))
		e.AsignarTropasRefuerzo(e.Jugadores[e.TurnoJugador])
	} else {
		// No se asignan nuevas tropas durante la fase de inicio
		e.Acciones = append(e.Acciones, NewAccionInicioTurno(e.Jugadores[e.TurnoJugador], 0, 0, 0))
	}
}

// AsignarTropasRefuerzo otorga un número de ejércitos al jugador que comienza un turno, dependiendo del número de territorios
// que ocupa. El número de ejércitos será la división entera del número de territorios ocupados entre 3.
// Cabe destacar que como mínimo, se otorgarán 3 ejércitos al principio de cada turno, independientemente de los
// territorios.
// Además, si el jugador controla por completo un continente, recibirá ejércitos extra. El número de ejercitos dependerá
// del continente:
//		- 2 ejércitos para Oceanía y América del Sur
//		- 3 ejércitos para África
//		- 5 ejércitos para América del Norte y Europa
//		- 7 ejércitos para Asia
func (e *EstadoPartida) AsignarTropasRefuerzo(jugador string) {
	regionesOcupadas := 0

	// Comprobamos el número de regiones que ocupa
	for _, region := range e.EstadoMapa {
		if region.Ocupante == jugador {
			regionesOcupadas++
		}
	}

	tropasObtenidas := 0
	if regionesOcupadas < 12 {
		tropasObtenidas = 3
	} else {
		tropasObtenidas = regionesOcupadas / 3
	}

	// Comprobamos si controla algún continente por completo
	continentesControlados := 0
	for _, c := range Continentes {
		puedeControlar := true
		for _, region := range c.Regiones {
			// Si alguna región del continente no es ocupada por el jugador, no lo puede controlar completamente
			if e.EstadoMapa[region].Ocupante != jugador {
				puedeControlar = false
				break
			}
		}

		// Si todas las regiones son ocupadas por el jugador, controla el continente
		if puedeControlar {
			continentesControlados++
			tropasObtenidas += c.Valor
		}
	}

	e.EstadosJugadores[jugador].Tropas += tropasObtenidas
	e.Acciones = append(e.Acciones, NewAccionInicioTurno(jugador, tropasObtenidas, regionesOcupadas, continentesControlados))
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

	if e.Fase != Refuerzo {
		return errors.New("Solo se puede reforzar durante la fase de refuerzo")
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
	carta, cartas, err := retirarPrimeraCarta(e.Cartas)
	if err != nil {
		// No quedan cartas en la baraja
		// Devolvemos los descartes a la baraja y barajamos
		e.Cartas = e.Descartes
		e.Descartes = nil
		barajarCartas(e.Cartas)
		carta, cartas, err = retirarPrimeraCarta(e.Cartas)

		// Si aun así no hay cartas, devuelve error
		if err != nil {
			return err
		}
	}
	e.Cartas = cartas
	estado.Cartas = append(estado.Cartas, carta)

	// Añadimos la acción
	e.Acciones = append(e.Acciones, NewAccionObtenerCarta(carta, jugador))
	return nil
}

// CambiarCartas permite al jugador cambiar un conjunto de 3 cartas por ejércitos.
// Los cambios válidos son los siguientes:
//		- 3 cartas del mismo tipo
//		- 2 cartas del mismo tipo más un comodín
//		- 3 cartas, una de cada tipo
// Los cambios se realizarán durante la fase de refuerzo, o en fase de ataque, si el jugador tiene más
// de 4 cartas tras derrotar a un rival.
// Si alguno de los territorios de las cartas cambiadas están ocupados por el jugador, recibirá tropas extra.
// El número de tropas recibidas dependerá del número de cambios totales:
// 		- En el primer cambio se recibirán 4 cartas
//		- Por cada cambio, se recibirán 2 cartas más que en el anterior
//		- En el sexto cambio se recibirán 15 cartas
// 		- A partir del sexto cambio, se recibirán 5 cartas más que en el cambio anterior
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
	carta1, cartas, _ := RetirarCartaPorID(ID_carta1, estado.Cartas)
	carta2, cartas, _ := RetirarCartaPorID(ID_carta2, cartas)
	carta3, cartas, _ := RetirarCartaPorID(ID_carta3, cartas)
	estado.Cartas = cartas

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

	regiones := obtenerRegionesCartas([]Carta{carta1, carta2, carta3})
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

// FinDeFase permite al jugador terminar la fase actual de su turno y pasar a la siguiente.
// Para ello, el jugador que quiera cambiar de fase deberá ser aquel que tenga el turno actual.
// Cada fase tendrá unas condiciones especiales para el cambio de turno.
// En el refuerzo, no podrá cambiar de fase si tiene más de 4 cartas o si le quedan tropas por asignar
// En el ataque, no podrá cambiar de fase si tiene más de 4 cartas o si tiene que ocupar un territorio y aún no lo ha hecho.
// En la fortificación podrá cambiar de fase (dándole el turno a otro jugador) libremente
func (e *EstadoPartida) FinDeFase(jugador string) error {
	if !e.esTurnoJugador(jugador) {
		return errors.New("Solo puedes cambiar de fase durante tu turno")
	}

	estadoJugador := e.EstadosJugadores[jugador]
	switch e.Fase {
	case Refuerzo:
		log.Println("NUM CARTAS", len(estadoJugador.Cartas))
		if len(estadoJugador.Cartas) > 4 {
			return errors.New("Estás obligado a cambiar cartas hasta tener menos de 5")
		}
		if estadoJugador.Tropas > 0 {
			return errors.New("Estás obligado a asignar todas tus tropas para cambiar de fase")
		}

		// Pasamos a fase de ataque
		e.Fase = Ataque
		e.Acciones = append(e.Acciones, NewAccionCambioFase(Ataque, jugador))
	case Ataque:
		if len(estadoJugador.Cartas) > 4 {
			return errors.New("Estás obligado a cambiar cartas hasta tener menos de 5")
		}

		for _, region := range e.EstadoMapa {
			if region.Ocupante == "" {
				return errors.New("No puedes finalizar la fase de ataque dejando territorios desocupados")
			}
		}

		// Pasamos a la fase de fortificación
		e.Fase = Fortificar
		e.Acciones = append(e.Acciones, NewAccionCambioFase(Fortificar, jugador))
	case Fortificar:
		// El jugador roba una carta si ha conquistado algún territorio y no ha robado ya
		if e.HaConquistado && !e.HaRecibidoCarta {
			err := e.RecibirCarta(jugador)
			if err != nil {
				return err
			}
		}

		// Pasamos el turno al siguiente jugador
		e.SiguienteJugador()
	}

	return nil
}

// ConsultarCartas devuelve un slice que contiene las cartas que posee el usuario "jugador"
func (e *EstadoPartida) ConsultarCartas(jugador string) []Carta {
	return e.EstadosJugadores[jugador].Cartas
}

func (e *EstadoPartida) esTurnoJugador(jugador string) bool {
	return e.Jugadores[e.TurnoJugador] == jugador
}

func (e *EstadoPartida) ObtenerJugadorTurno() string {
	return e.Jugadores[e.TurnoJugador]
}
