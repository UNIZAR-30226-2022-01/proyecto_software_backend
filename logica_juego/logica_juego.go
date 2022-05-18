// Package logica_juego define estructuras y constantes para los diferentes elementos del juego,
// así como la lista de acciones y regiones a ser interpretada por los clientes.
package logica_juego

import (
	"errors"
	"gonum.org/v1/gonum/graph/simple"
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

	// Vector de booleanos con una entrada por cada jugador
	// Si el jugador sigue en la partida, su entrada valdrá true
	// Si ha sido derrotado, será false
	JugadoresActivos []bool

	Fase Fase

	EstadoMapa map[NumRegion]*EstadoRegion

	// Baraja
	Cartas     []Carta
	Descartes  []Carta
	NumCambios int

	// Recibir carta
	HaConquistado   bool // True si ha conquistado algún territorio en el turno
	HaRecibidoCarta bool // True si ya ha robado carta, para evitar más de un robo
	HaFortificado   bool // True si ya ha fortificado en el turno

	// Información sobre el último ataque
	RegionUltimoAtaque         NumRegion // Región desde la que se inició el último ataque
	DadosUltimoAtaque          int       // Número de dados que lanzo el atacante en el último ataque
	TropasPerdidasUltimoAtaque int       // Número de tropas que perdió el atacante en el último ataque
	HayTerritorioDesocupado    bool      // True si hay algún territorio sin ocupar
	UltimoDefensor             string    // Nombre del jugador defensor en el último ataque

	// Flag de partida terminada, pendiente por ser tener su estado consultado por todos los jugadores
	// previo a su eliminación
	Terminada                      bool
	JugadoresRestantesPorConsultar []string

	// Timestamp de la última acción realizada en el juego por el usuario del turno actual, para
	// tratar la eliminación de usuarios inactivos
	UltimaAccion time.Time

	// Flag de si se ha enviado una alerta por inactividad al jugador actual
	AlertaEnviada bool
}

func CrearEstadoPartida(jugadores []string) (e EstadoPartida) {
	jugadoresActivos := make([]bool, len(jugadores))
	for i := range jugadoresActivos {
		jugadoresActivos[i] = true
	}
	e = EstadoPartida{
		Acciones:                       make([]interface{}, 0),
		Jugadores:                      crearSliceJugadores(jugadores),
		EstadosJugadores:               crearMapaEstadosJugadores(jugadores),
		TurnoJugador:                   (LanzarDados()) % len(jugadores), // Primer jugador aleatorio
		Fase:                           Inicio,
		EstadoMapa:                     crearEstadoMapa(),
		Cartas:                         crearBaraja(),
		Descartes:                      []Carta{},
		NumCambios:                     0,
		HaConquistado:                  false,
		HaRecibidoCarta:                false,
		HaFortificado:                  false,
		DadosUltimoAtaque:              0,
		TropasPerdidasUltimoAtaque:     0,
		HayTerritorioDesocupado:        false,
		JugadoresActivos:               jugadoresActivos,
		Terminada:                      false,
		JugadoresRestantesPorConsultar: crearSliceJugadores(jugadores),
		UltimaAccion:                   time.Now(),
		AlertaEnviada:                  false,
	}

	return e
}

// ExpulsarJugadorActual contabiliza el jugador del turno actual (derrotado), añade una acción de expulsión y
// pasa al siguiente jugador si es posible
func (e *EstadoPartida) ExpulsarJugadorActual() {
	jugador := e.Jugadores[e.TurnoJugador]

	e.JugadoresActivos[e.obtenerTurnoJugador(jugador)] = false

	e.Acciones = append(e.Acciones, NewAccionJugadorExpulsado(jugador))

	jugadoresActivos := 0
	for _, act := range e.JugadoresActivos {
		if act {
			jugadoresActivos += 1
		}
	}

	if jugadoresActivos > 0 {
		e.SiguienteJugador()
	} else {
		// Todos los jugadores han sido derrotados o expulsados, se para la Goroutine y espera
		// a que consulten el estado
		e.Terminada = true
		e.JugadoresRestantesPorConsultar = nil
	}
}

// ExpulsarJugador contabiliza un jugador dado como expulsado (derrotado), añade una acción de expulsión y
// pasa al siguiente jugador si era el actual y es posible
func (e *EstadoPartida) ExpulsarJugador(expulsado string) {
	e.JugadoresActivos[e.obtenerTurnoJugador(expulsado)] = false

	e.Acciones = append(e.Acciones, NewAccionJugadorExpulsado(expulsado))

	jugadoresActivos := 0
	for _, act := range e.JugadoresActivos {
		if act {
			jugadoresActivos += 1
		}
	}

	if jugadoresActivos > 0 && e.Jugadores[e.TurnoJugador] == expulsado { // Era el jugador del turno actual
		e.SiguienteJugador()
	} else if jugadoresActivos == 0 {
		// Todos los jugadores han sido derrotados o expulsados, se para la Goroutine y espera
		// a que consulten el estado
		e.Terminada = true
		e.JugadoresRestantesPorConsultar = nil
	}
}

// TerminadaPorExpulsiones devuelve true si la partida ha terminado por tener todos sus jugadores expulsados, false en otro caso
func (e *EstadoPartida) TerminadaPorExpulsiones() bool {
	return e.Terminada && len(e.JugadoresRestantesPorConsultar) == 0
}

// HaSidoEliminado devuelve true si el usuario ha sido eliminado por otro jugador,
// false en otro caso
// El nombre de jugador indicado debe existir en la partida
func (e *EstadoPartida) HaSidoEliminado(jugador string) bool {
	return e.ContarTerritoriosOcupados(jugador) == 0
}

// HaSidoExpulsado devuelve true si el usuario ha sido eliminado por inactividad,
// false en otro caso
// El nombre de jugador indicado debe existir en la partida
func (e *EstadoPartida) HaSidoExpulsado(jugador string) bool {
	return !e.JugadoresActivos[e.obtenerTurnoJugador(jugador)] && e.ContarTerritoriosOcupados(jugador) != 0
}

// SiguienteJugadorSinAccion cambia el turno a otro jugador sin emitir ninguna acción.
// Usar únicamente al rellenar regiones.
func (e *EstadoPartida) SiguienteJugadorSinAccion() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)
}

// SiguienteJugador cambia el turno a otro jugador, emitiendo la acción correspondiente.
// TODO SiguienteJugador no debería ser pública, lo es por necesidad dentro del test
func (e *EstadoPartida) SiguienteJugador() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)

	for !e.JugadoresActivos[e.TurnoJugador] {
		// Saltamos al jugador en caso de que haya sido derrotado
		e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)
	}

	// En el nuevo turno no se habrá recibido carta, conquistado ni fortificado
	e.HaRecibidoCarta = false
	e.HaConquistado = false
	e.HaFortificado = false

	// Reiniciamos el estado del último ataque
	e.DadosUltimoAtaque = 0
	e.TropasPerdidasUltimoAtaque = 0
	e.HayTerritorioDesocupado = false

	if e.Fase != Inicio {
		e.Acciones = append(e.Acciones, NewAccionCambioFase(Refuerzo, e.Jugadores[e.TurnoJugador]))
		// Se introduce una acción de nuevo turno al asignar tropas de refuerzo
		e.AsignarTropasRefuerzo(e.Jugadores[e.TurnoJugador])
		// Pasamos a la fase de refuerzo
		e.Fase = Refuerzo
	} else {
		// No se asignan nuevas tropas durante la fase de inicio
		e.Acciones = append(e.Acciones, NewAccionCambioFase(Inicio, e.Jugadores[e.TurnoJugador]))
		e.Acciones = append(e.Acciones, NewAccionInicioTurno(e.Jugadores[e.TurnoJugador], 0, 0, 0))
	}

	// Refresca el timestamp
	e.UltimaAccion = time.Now()

	e.AlertaEnviada = false
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
// TODO: RellenarRegiones más elaborado (pseudo-random teniendo en cuenta adyacencias, recorridos por el grafo, etc.)
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

			// La añade al subgrafo del jugador
			//e.EstadosJugadores[e.Jugadores[e.TurnoJugador]].AñadirASubGrafo(i)

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
		return errors.New("Solo puedes reforzar durante tu turno, el jugador en este turno es " + e.ObtenerJugadorTurno())
	}

	if e.Fase == Ataque || e.Fase == Fortificar {
		return errors.New("Solo se puede reforzar durante la fase de refuerzo o durante el inicio de la partida")
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

// FinDeFase permite al jugador terminar la fase actual de su turno y pasar a la siguiente.
// Para ello, el jugador que quiera cambiar de fase deberá ser aquel que tenga el turno actual.
// Cada fase tendrá unas condiciones especiales para el cambio de turno.
// En el refuerzo, no podrá cambiar de fase si tiene más de 4 cartas o si le quedan tropas por asignar
// En el ataque, no podrá cambiar de fase si tiene más de 4 cartas o si tiene que ocupar un territorio y aún no lo ha hecho.
// En la fortificación podrá cambiar de fase (dándole el turno a otro jugador) libremente
func (e *EstadoPartida) FinDeFase(jugador string) error {
	if !e.esTurnoJugador(jugador) {
		return errors.New("Solo puedes cambiar de fase durante tu turno, el jugador en este turno es " + e.ObtenerJugadorTurno())
	}

	estadoJugador := e.EstadosJugadores[jugador]
	switch e.Fase {
	case Inicio:
		if len(estadoJugador.Cartas) > 4 {
			return errors.New("Estás obligado a cambiar cartas hasta tener menos de 5")
		}
		if estadoJugador.Tropas > 0 {
			return errors.New("Estás obligado a asignar todas tus tropas para cambiar de fase, te quedan " + strconv.Itoa(estadoJugador.Tropas) + " tropas")
		}
		// Se comprueba si la fase ha finalizado (no le quedan tropas a ningún jugador)
		todosSinTropas := true
		for _, jugador := range e.Jugadores {
			estadoJugador := e.EstadosJugadores[jugador]

			// Si aún tiene tropas y no ha sido expulsado ni eliminado, se sigue en fase de inicio
			if estadoJugador.Tropas != 0 && !e.HaSidoExpulsado(jugador) && !e.HaSidoEliminado(jugador) {
				todosSinTropas = false
				break
			}
		}

		if todosSinTropas { // Se pasa al siguiente jugador, que empezará en fase de Refuerzo
			e.Fase = Fortificar // Se finge que se ha pasado desde fortificar (o cualquier otra excepto Inicio)
		} //Si no, se mantiene en la fase de inicio para el siguiente jugador

		e.SiguienteJugador()

	case Refuerzo:
		if len(estadoJugador.Cartas) > 4 {
			return errors.New("Estás obligado a cambiar cartas hasta tener menos de 5")
		}
		if estadoJugador.Tropas > 0 {
			return errors.New("Estás obligado a asignar todas tus tropas para cambiar de fase, te quedan " + strconv.Itoa(estadoJugador.Tropas) + " tropas")
		}

		// Pasamos a fase de ataque
		e.Fase = Ataque
		e.Acciones = append(e.Acciones, NewAccionCambioFase(Ataque, jugador))

	case Ataque:
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

	// Refresca el timestamp
	e.UltimaAccion = time.Now()

	e.AlertaEnviada = false

	return nil
}

func (e *EstadoPartida) esTurnoJugador(jugador string) bool {
	return e.Jugadores[e.TurnoJugador] == jugador
}

func (e *EstadoPartida) ObtenerJugadorTurno() string {
	return e.Jugadores[e.TurnoJugador]
}

// obtenerTurnoJugador devuelve un entero indicando el turno correspondiente al jugador determinado
// Es decir, devuelve su posición dentro del vector de jugadores
func (e *EstadoPartida) obtenerTurnoJugador(jugador string) int {
	for i, j := range e.Jugadores {
		if j == jugador {
			return i
		}
	}

	return -1
}

func (e *EstadoPartida) FortificarTerritorio(origen int, destino int, tropas int, jugador string) error {
	if !e.esTurnoJugador(jugador) {
		return errors.New("Solo puedes fortificar durante tu turno, el jugador en este turno es " + e.ObtenerJugadorTurno())
	} else if e.Fase != Fortificar {
		return errors.New("Solo se puede fortificar durante la fase de fortificación")
	}

	if e.HaFortificado {
		return errors.New("No puedes fortificar más de una vez por turno")
	}
	// Comprobar pertenencia de territorio de origen y número de tropas válido en el territorio
	regionOrigen, existe := e.EstadoMapa[NumRegion(origen)]
	if !existe {
		return errors.New("El territorio de origen indicado no existe")
	} else if regionOrigen.Ocupante != jugador {
		return errors.New("No eres el ocupante del territorio de origen, el ocupante es " + regionOrigen.Ocupante)
	} else if regionOrigen.NumTropas == 1 {
		return errors.New("El número de tropas en el territorio de origen debe ser mayor que 1")
	} else if regionOrigen.NumTropas <= tropas {
		return errors.New("El número de tropas en el territorio de origen debe ser mayor que el número de tropas de fortificación")
	}

	// Comprobar pertenencia de territorio destino
	regionDestino, existe := e.EstadoMapa[NumRegion(destino)]
	if !existe {
		return errors.New("El territorio de destino indicado no existe")
	} else if regionDestino.Ocupante != jugador {
		return errors.New("No eres el ocupante del territorio de destino, el ocupante es " + regionDestino.Ocupante)
	} else if regionDestino == regionOrigen {
		return errors.New("Las regiones origen y destino deben ser diferentes")
	}

	// Comprobar existencia de un camino entre ambas regiones, que cruce exclusivamente por territorios controlados
	existe = e.existeCaminoEntreRegiones(NumRegion(origen), NumRegion(destino), jugador)

	if !existe {
		return errors.New("No existe un camino con territorios controlados entre ambas regiones")
	} else {
		// Realizar la fortificación
		e.EstadoMapa[NumRegion(origen)].NumTropas -= tropas
		e.EstadoMapa[NumRegion(destino)].NumTropas += tropas

		// Indicamos que ya se ha fortificado en ese turno
		e.HaFortificado = true

		// Añadir una nueva acción
		e.Acciones = append(e.Acciones,
			NewAccionFortificar(
				NumRegion(origen),
				NumRegion(destino),
				e.EstadoMapa[NumRegion(origen)].NumTropas,
				e.EstadoMapa[NumRegion(destino)].NumTropas,
				jugador))
	}

	return nil
}

// Implementación del algoritmo BFS para buscar e indicar si existe un camino
// entre dos regiones de un sub-grafo de un jugador, dados sus identificadores
func (e *EstadoPartida) existeCaminoEntreRegiones(origen NumRegion, destino NumRegion, jugador string) bool {
	//region := e.EstadoMapa[NumRegion(origen)]

	// No se puede usar un canal, por ejemplo, porque tienen un límite fijo de espacio
	var frontera []NumRegion
	frontera = append(frontera, origen)

	var explorados []NumRegion

	for {
		if len(frontera) == 0 {
			return false
		}

		// "pop"
		idRegion := frontera[0]
		frontera = frontera[1:]

		explorados = append(explorados, idRegion)

		nodo := simple.Node(idRegion)
		nodos := e.ObtenerSubgrafoRegiones(jugador).From(nodo.ID())

		// https://pkg.go.dev/gonum.org/v1/gonum/graph@v0.11.0#Iterator
		for nodos.Next() == true { // La siguiente llamada no devolverá un nil
			hijo := nodos.Node()

			if !regionEnCola(explorados, NumRegion(hijo.ID())) && !regionEnCola(frontera, NumRegion(hijo.ID())) {
				if NumRegion(hijo.ID()) == destino {
					return true
				}
				frontera = append(frontera, NumRegion(hijo.ID()))
			}
		}
	}
}

func regionEnCola(cola []NumRegion, region NumRegion) bool {
	for _, regionCola := range cola {
		if region == regionCola {
			return true
		}
	}

	return false
}

// ObtenerSubgrafoRegiones devuelve un subgrafo de GrafoMapa con únicamente las regiones/nodos controladas por el jugador
// TODO: La librería de grafos no exporta sus campos y no se puede serializar junto al resto del estado,
// TODO: por lo que se tienen que recrear o almacenar cacheados en otra estructura que no se serialice
// TODO: Valorar el coste de almacenamiento vs. recreación en cada llamada a fortificar (el coste es de
// TODO: NumRegiones*factor_ramificacion_maximo (máximo número de aristas en un nodo) iteraciones, exactamente)
func (e *EstadoPartida) ObtenerSubgrafoRegiones(jugador string) *simple.UndirectedGraph {
	subgrafo := simple.NewUndirectedGraph()

	for i := Eastern_australia; i <= Alberta; i++ {
		if e.EstadoMapa[i].Ocupante == jugador {
			AñadirASubGrafo(i, subgrafo)
		}
	}

	return subgrafo
}

// AñadirASubGrafo añade una región dada al subgrafo del jugador, conectándola con el resto
// de regiones del subgrafo que se encontraran conectadas con ella en el grafo
// del mapa completo
func AñadirASubGrafo(region NumRegion, subgrafo *simple.UndirectedGraph) {
	// Añade la región al grafo
	regionNueva := simple.Node(region)
	subgrafo.AddNode(regionNueva) // TODO: no tiene tratamiento de errores, hace panic, encerrar en búsqueda

	// Para cada región alcanzable desde regionNueva
	regionesAlcanzables := GrafoMapa.From(regionNueva.ID())
	for regionesAlcanzables.Next() == true { // La siguiente llamada no devolverá un nil
		regionAlcanzable := regionesAlcanzables.Node()

		if subgrafo.Node(regionAlcanzable.ID()) != nil { // Si la región alcanzable está en el subgrafo, se conectan
			subgrafo.SetEdge(GrafoMapa.NewEdge(regionNueva, regionAlcanzable))
		}
	}
}

// Elimina una región dada del subgrafo del jugador, desconectándola también del
// resto de regiones del subgrafo que se encontraran conectadas con ella en el
// grafo del mapa completo
func (e *EstadoJugador) eliminarDeSubgrafo(region NumRegion, subgrafo *simple.UndirectedGraph) {
	subgrafo.RemoveNode(int64(region)) // TODO: no tiene tratamiento de errores, es una no-op si no existe, encerrar en búsqueda
}

// EnviarMensaje encola una acción en la lista de acciones de la partida, que representa el envío de un mensaje por parte
// de "jugador", con el contenido "mensaje"
func (e *EstadoPartida) EnviarMensaje(jugador, mensaje string) {
	e.Acciones = append(e.Acciones, NewAccionMensaje(jugador, mensaje))
}
