package vo

import (
	"backend/logica_juego"
	"errors"
	"log"
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
	Jugador             string
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

func CrearEstadoPartida(jugadores []Usuario) (e EstadoPartida) {
	e = EstadoPartida{
		Acciones:         make([]interface{}, 0),
		Jugadores:        crearSliceJugadores(jugadores),
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

func crearSliceJugadores(jugadores []Usuario) (slice []string) {
	for _, jugador := range jugadores {
		slice = append(slice, jugador.NombreUsuario)
	}

	return slice
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

	log.Println("mapa estados jugadores:", mapa)

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
			e.Acciones = append(e.Acciones, RecibirRegion{
				IDAccion:             0, // TODO: enum
				Region:               i,
				TropasRestantes:      e.EstadosJugadores[e.Jugadores[e.TurnoJugador]].Tropas,
				TerritoriosRestantes: NUM_REGIONES - regionesAsignadas,
				Jugador:              e.Jugadores[e.TurnoJugador],
			})

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

	if estado.Tropas-numTropas < 0 {
		return errors.New("No tienes tropas suficientes para reforzar un territorio, tropas restantes: " + strconv.Itoa(estado.Tropas))
	} else {
		estado.Tropas = estado.Tropas - numTropas
		region.NumTropas = region.NumTropas + numTropas

		// Añadir una nueva acción
		e.Acciones = append(e.Acciones, AccionReforzar{
			IDAccion:            0, // TODO: enum
			Jugador:             jugador,
			TerritorioReforzado: NumRegion(idTerritorio),
			TropasRefuerzo:      numTropas,
		})

		return nil
	}
}

func (e *EstadoPartida) esTurnoJugador(jugador string) bool {
	return e.Jugadores[e.TurnoJugador] == jugador
}

func (e *EstadoPartida) ObtenerJugadorTurno() string {
	return e.Jugadores[e.TurnoJugador]
}
