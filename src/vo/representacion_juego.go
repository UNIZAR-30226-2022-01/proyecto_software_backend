package vo

import (
	"backend/logica_juego"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type Fase int
type TipoTropa int
type NumRegion int

const (
	Inicio Fase = iota // Repartir regiones
	Refuerzo
	Ataque
	Fortificar
)

const (
	Infanteria TipoTropa = iota
	Caballeria
	Artilleria
)

const (
	NUM_REGIONES = 42
)

const (
	Eastern_australia NumRegion = iota
	Indonesia
	New_guinea
	Alaska
	Ontario
	Northwest_territory
	Venezuela
	Madagascar
	North_africa
	Greenland
	Iceland
	Great_britain
	Scandinavia
	Japan
	Yakursk
	Kamchatka
	Siberia
	Ural
	Afghanistan
	Middle_east
	India
	Siam
	China
	Mongolia
	Irkutsk
	Ukraine
	Southern_europe
	Western_europe
	Northern_europe
	Egypt
	East_africa
	Congo
	South_africa
	Brazil
	Argentina
	Eastern_united_states
	Western_united_states
	Quebec
	Central_america
	Peru
	Western_australia
	Alberta
)

func (nr NumRegion) String() string {
	// Lo sentimos
	return []string{"eastern_australia", "indonesia",
		"new_guinea", "alaska", "ontario", "northwest_territory",
		"venezuela", "madagascar", "north_africa", "greenland",
		"iceland", "great_britain", "scandinavia", "japan", "yakursk",
		"kamchatka", "siberia", "ural", "afghanistan", "middle_east",
		"india", "siam", "china", "mongolia", "irkutsk", "ukraine",
		"southern_europe", "western_europe", "northern_europe", "egypt",
		"east_africa", "congo", "south_africa", "brazil", "argentina",
		"eastern_united_states", "western_united_states", "quebec",
		"central_america", "peru", "western_australia", "alberta"}[nr]
}

type EstadoRegion struct {
	Ocupante  string
	NumTropas int
}

type Carta struct {
	Tipo      TipoTropa
	Region    NumRegion
	EsComodin bool
}

/*
Lista de acciones:
*/

type RecibirRegion struct {
	IDAccion             int
	Region               NumRegion
	TropasRestantes      int
	TerritoriosRestantes int
	Jugador              string
}

type AccionCambioTurno struct {
	IDAccion int
	Turno    int
	Jugador  string
}

type AccionCambioFase struct {
	IDAccion int
	Fase     Fase
	Jugador  string
}

type AccionInicioTurno struct {
	IDAccion                 int
	Jugador                  string
	TropasObtenidas          int
	RazonNumeroTerritorios   int
	RazonContinentesOcupados int
}

type AccionCambioCartas struct {
	IDAccion                    int
	NumConjuntosCambiados       int
	NumTropasObtenidas          int
	BonificacionObtenida        bool
	RegionQueOtorgaBonificacion NumRegion
	ObligadoAHacerCambios       bool
}

type AccionReforzar struct {
	IDAccion            int
	TerritorioReforzado NumRegion
	TropasRefuerzo      int
}

type AccionAtaque struct {
	IDAccion               int
	Origen                 NumRegion
	Destion                NumRegion
	TropasPerdidasAtacante NumRegion
	TropasPerdidasDefensor NumRegion
	NumDadosAtaque         int
	JugadorAtacante        string
	JugadorDefensor        string
}

type AccionOcupar struct {
	IDAccion        int
	Origen          NumRegion
	Destion         NumRegion
	TropasOrigen    NumRegion
	TropasDestino   NumRegion
	JugadorOcupante string
	JugadorOcupado  string
}

type AccionFortificar struct {
	IDAccion      int
	Origen        NumRegion
	Destino       NumRegion
	TropasOrigen  int
	TropasDestino int
	Jugador       string
}

type AccionObtenerCarta struct {
	IDAccion int
	Carta    Carta
	Jugador  string
}

type EstadoJugador struct {
	Cartas            []Carta
	UltimoIndiceLeido int
	Tropas            int
}

type EstadoPartida struct {
	Acciones         []interface{}             // Acciones realizadas durante la partida
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

func CrearEstadoPartida(jugadores []Usuario) (e *EstadoPartida) {
	*e = EstadoPartida{
		Acciones:         make([]interface{}, 0),
		EstadosJugadores: crearMapaEstadosJugadores(jugadores),
		TurnoJugador:     logica_juego.LanzarDados(), // Primer jugador aleatorio
		Fase:             Inicio,
		NumeroTurno:      0,
		EstadoMapa:       crearEstadoMapa(),
		Cartas:           crearBaraja(),
		NumCambios:       0,
	}

	return e
}

func crearEstadoMapa() (mapa map[NumRegion]*EstadoRegion) {
	mapa = make(map[NumRegion]*EstadoRegion)
	for i := Eastern_australia; i <= Alberta; i++ {
		mapa[i] = &EstadoRegion{Ocupante: "", NumTropas: 0}
	}

	return mapa
}

func crearMapaEstadosJugadores(jugadores []Usuario) (mapa map[string]*EstadoJugador) {
	mapa = make(map[string]*EstadoJugador, len(jugadores))

	// Tropas iniciales, según el número de jugadores
	numTropas := 0
	if len(jugadores) == 3 {
		numTropas = 35
	} else if len(jugadores) == 4 {
		numTropas = 30
	} else if len(jugadores) == 5 {
		numTropas = 25
	} else if len(jugadores) == 6 {
		numTropas = 20
	}

	for _, j := range jugadores {
		mapa[j.NombreUsuario] = &EstadoJugador{}

		mapa[j.NombreUsuario].Cartas = []Carta{}
		mapa[j.NombreUsuario].UltimoIndiceLeido = -1
		mapa[j.NombreUsuario].Tropas = numTropas
	}

	return mapa
}

func crearBaraja() (cartas []Carta) {
	numTiposTropa := 0
	for i := Eastern_australia; i <= Alberta; i++ {
		var carta Carta
		var tipo TipoTropa
		if numTiposTropa < 18 {
			tipo = Infanteria
		} else if numTiposTropa < 36 {
			tipo = Caballeria
		} else {
			tipo = Artilleria
		}

		carta = Carta{
			Tipo:      tipo,
			Region:    i,
			EsComodin: false,
		}

		numTiposTropa = numTiposTropa + 1

		cartas = append(cartas, carta)
	}

	cartas = append(cartas, Carta{EsComodin: true})
	cartas = append(cartas, Carta{EsComodin: true})

	// Se baraja aleatoriamente
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cartas), func(i, j int) { cartas[i], cartas[j] = cartas[j], cartas[i] })

	return cartas
}

func (e *EstadoPartida) SiguienteJugador() {
	e.TurnoJugador = (e.TurnoJugador + 1) % len(e.EstadosJugadores)
}

// TODO: Más elaborado (pseudo-random, con recorridos por el grafo, etc.)
func (e *EstadoPartida) RellenarRegiones() {
	// Extrae los nombres de jugador del mapa
	jugadores := make([]string, len(e.EstadosJugadores))
	for jugador := range e.EstadosJugadores {
		jugadores = append(jugadores, jugador)
	}

	regionesAsignadas := 0
	for i := Eastern_australia; i <= Alberta; i++ {
		if e.EstadosJugadores[jugadores[e.TurnoJugador]].Tropas >= 1 {
			e.EstadoMapa[i].Ocupante = jugadores[e.TurnoJugador]
			e.EstadoMapa[i].NumTropas = 1

			e.EstadosJugadores[jugadores[e.TurnoJugador]].Tropas = e.EstadosJugadores[jugadores[e.TurnoJugador]].Tropas - 1

			regionesAsignadas = regionesAsignadas + 1

			// Añadir una nueva acción
			e.Acciones = append(e.Acciones, RecibirRegion{
				IDAccion:             0,
				Region:               i,
				TropasRestantes:      e.EstadosJugadores[jugadores[e.TurnoJugador]].Tropas,
				TerritoriosRestantes: NUM_REGIONES - regionesAsignadas,
				Jugador:              jugadores[e.TurnoJugador],
			})
		} else {
			// Repite la iteración para el siguiente jugador
			e.SiguienteJugador()
			i = i - 1
		}
	}
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
	} else if e.esTurnoJugador(jugador) {
		return errors.New("Se ha solicitado una acción fuera de turno, el jugador en este turno es " + e.obtenerJugadorTurno())
	}

	if estado.Tropas-numTropas < 0 {
		return errors.New("No tienes tropas suficientes para reforzar un territorio, tropas restantes: " + strconv.Itoa(estado.Tropas))
	} else {
		estado.Tropas = estado.Tropas - 1
		region.NumTropas = region.NumTropas + numTropas

		// Añadir una nueva acción
		e.Acciones = append(e.Acciones, AccionReforzar{
			IDAccion:            0,
			TerritorioReforzado: NumRegion(idTerritorio),
			TropasRefuerzo:      numTropas,
		})

		return nil
	}
}

func (e *EstadoPartida) esTurnoJugador(jugador string) bool {
	// Extrae los nombres de jugador del mapa
	jugadores := make([]string, len(e.EstadosJugadores))
	for jugador := range e.EstadosJugadores {
		jugadores = append(jugadores, jugador)
	}

	return jugadores[e.TurnoJugador] == jugador
}

func (e *EstadoPartida) obtenerJugadorTurno() string {
	// Extrae los nombres de jugador del mapa
	jugadores := make([]string, len(e.EstadosJugadores))
	for jugador := range e.EstadosJugadores {
		jugadores = append(jugadores, jugador)
	}

	return jugadores[e.TurnoJugador]
}
