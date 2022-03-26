package vo

import (
	"math/rand"
	"net/http"
	"time"
)

type Fase int
type TipoTropa int
type NumRegion int

const (
	Refuerzo Fase = iota
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

// Usuario es un objeto de usuario equivalente al del modelo de base de datos.
type Usuario struct {
	Email           string
	NombreUsuario   string
	PasswordHash    string
	Biografia       string
	CookieSesion    http.Cookie
	Puntos          int
	PartidasGanadas int
	PartidasTotales int
	ID_dado         int
	ID_ficha        int
}

type Partida struct {
	IdPartida          int
	EsPublica          bool
	PasswordHash       string
	EnCurso            bool
	MaxNumeroJugadores int
	Jugadores          []Usuario

	// TODO: representar estado del juego y chat de la partida
	Mensajes []Mensaje
	Estado   EstadoPartida
}

type Mensaje struct {
	// TODO
}

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

//func (act AccionCambioTurno) A() {}

type EstadoPartida struct {
	Acciones    []interface{}
	Jugador     string
	Fase        Fase
	NumeroTurno int

	EstadoMapa map[NumRegion]EstadoRegion

	// Baraja
	Cartas          []Carta
	CartasJugadores map[string][]Carta
	NumCambios      int

	// ...
}

func (p *Partida) InicializarAcciones() {
	p.Estado.Acciones = make([]interface{}, 0)

	//p.Estado.acciones = append(p.Estado.acciones, AccionCambioTurno{29278927289728})
	//p.Estado.acciones = append(p.Estado.acciones, AccionCambioTurno2{"adsasdadsadsads"})
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

func (p *Partida) crearMapaCartasJugadores() {
	p.Estado.CartasJugadores = make(map[string][]Carta, p.MaxNumeroJugadores)

	for _, j := range p.Jugadores {
		p.Estado.CartasJugadores[j.NombreUsuario] = []Carta{}
	}
}

func (p *Partida) CrearEstadoPartida() {
	p.Estado = EstadoPartida{
		Jugador:         "",
		Fase:            Refuerzo,
		NumeroTurno:     0,
		EstadoMapa:      crearEstadoMapa(),
		Cartas:          crearBaraja(),
		CartasJugadores: nil,
		NumCambios:      0,
	}

	p.crearMapaCartasJugadores()
}

///////////////////////////////////////////
// Structs de respuesta a serializar a JSON
///////////////////////////////////////////

type ElementoListaPartidas struct {
	IdPartida          int
	EsPublica          bool
	NumeroJugadores    int
	MaxNumeroJugadores int
	AmigosPresentes    []string
	NumAmigosPresentes int
}

// ContarAmigos devuelve cuántos amigos de un usuario están en una partida dada
func ContarAmigos(amigos []Usuario, partida Partida) (num int) {
	for _, amigo := range amigos {
		// Como máximo hay 6 jugadores en la partida, así que
		// la complejidad la dicta el número de amigos del usuario
		for _, jugador := range partida.Jugadores {
			if amigo.NombreUsuario == jugador.NombreUsuario {
				num = num + 1
			}
		}
	}

	return num
}
