package vo

import (
	"math/rand"
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
	IDAccion int
	Region   NumRegion
	Jugador  string
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

type EstadoPartida struct {
	Acciones     []interface{}  // Acciones realizadas durante la partida
	Jugadores    map[string]int // Mapa de nombres de los jugadores en la partida y su último índice no leído de la lista de acciones
	TurnoJugador int            // Índice de la lista que corresponde a qué jugador le toca

	Fase        Fase
	NumeroTurno int

	EstadoMapa map[NumRegion]EstadoRegion

	// Baraja
	Cartas          []Carta
	CartasJugadores map[string][]Carta
	NumCambios      int

	// ...
}

func crearEstadoMapa() (mapa map[NumRegion]EstadoRegion) {
	mapa = make(map[NumRegion]EstadoRegion)
	for i := Eastern_australia; i <= Alberta; i++ {
		mapa[i] = EstadoRegion{Ocupante: "", NumTropas: 0}
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

func crearMapaCartasJugadores(e *EstadoPartida, jugadores []Usuario) {
	e.CartasJugadores = make(map[string][]Carta, len(jugadores))

	for _, j := range jugadores {
		e.CartasJugadores[j.NombreUsuario] = []Carta{}
	}
}

func crearMapaIndicesJugadores(e *EstadoPartida, jugadores []Usuario) {
	e.Jugadores = make(map[string]int, len(jugadores))

	for _, j := range jugadores {
		e.Jugadores[j.NombreUsuario] = -1 // Empezarán leyendo el índice 0
	}
}

func CrearEstadoPartida(jugadores []Usuario) (e *EstadoPartida) {
	*e = EstadoPartida{
		Acciones:        make([]interface{}, 0),
		Jugadores:       nil,
		TurnoJugador:    0,
		Fase:            Inicio,
		NumeroTurno:     0,
		EstadoMapa:      crearEstadoMapa(),
		Cartas:          crearBaraja(),
		CartasJugadores: nil,
		NumCambios:      0,
	}

	crearMapaCartasJugadores(e, jugadores)
	crearMapaIndicesJugadores(e, jugadores)

	return e
}
