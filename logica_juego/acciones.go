package logica_juego

const (
	IDAccionRecibirRegion Fase = iota
	IDAccionCambioTurno
	IDAccionCambioFase
	IDAccionInicioTurno
	IDAccionCambioCartas
	IDAccionReforzar
	IDAccionAtaque
	IDAccionOcupar
	IDAccionFortificar
	IDAccionObtenerCarta
)

type AccionRecibirRegion struct {
	IDAccion             int
	Region               NumRegion
	TropasRestantes      int
	TerritoriosRestantes int
	Jugador              string
}

func NewAccionRecibirRegion(region NumRegion, tropasRestantes int, territoriosRestantes int, jugador string) AccionRecibirRegion {
	return AccionRecibirRegion{
		IDAccion:             int(IDAccionRecibirRegion),
		Region:               region,
		TropasRestantes:      tropasRestantes,
		TerritoriosRestantes: territoriosRestantes,
		Jugador:              jugador,
	}
}

type AccionCambioTurno struct {
	IDAccion int
	Jugador  string
}

func NewAccionCambioTurno(jugador string) AccionCambioTurno {
	return AccionCambioTurno{
		IDAccion: int(IDAccionCambioTurno),
		Jugador:  jugador}
}

type AccionCambioFase struct {
	IDAccion int
	Fase     Fase
	Jugador  string
}

func NewAccionCambioFase(fase Fase, jugador string) AccionCambioFase {
	return AccionCambioFase{
		IDAccion: int(IDAccionCambioFase),
		Fase:     fase,
		Jugador:  jugador}
}

type AccionInicioTurno struct {
	IDAccion                 int
	Jugador                  string
	TropasObtenidas          int
	RazonNumeroTerritorios   int
	RazonContinentesOcupados int
}

func NewAccionInicioTurno(jugador string, tropasObtenidas int, razonNumeroTerritorios int, razonContinentesOcupados int) AccionInicioTurno {
	return AccionInicioTurno{
		IDAccion:                 int(IDAccionInicioTurno),
		Jugador:                  jugador,
		TropasObtenidas:          tropasObtenidas,
		RazonNumeroTerritorios:   razonNumeroTerritorios,
		RazonContinentesOcupados: razonContinentesOcupados}
}

type AccionCambioCartas struct {
	IDAccion                    int
	NumConjuntosCambiados       int
	NumTropasObtenidas          int
	BonificacionObtenida        bool
	RegionQueOtorgaBonificacion NumRegion
	ObligadoAHacerCambios       bool
}

func NewAccionCambioCartas(numConjuntosCambiados int, numTropasObtenidas int, bonificacionObtenida bool, regionQueOtorgaBonificacion NumRegion, obligadoAHacerCambios bool) AccionCambioCartas {
	return AccionCambioCartas{
		IDAccion:                    int(IDAccionCambioCartas),
		NumConjuntosCambiados:       numConjuntosCambiados,
		NumTropasObtenidas:          numTropasObtenidas,
		BonificacionObtenida:        bonificacionObtenida,
		RegionQueOtorgaBonificacion: regionQueOtorgaBonificacion,
		ObligadoAHacerCambios:       obligadoAHacerCambios}
}

type AccionReforzar struct {
	IDAccion            int
	Jugador             string
	TerritorioReforzado NumRegion
	TropasRefuerzo      int
}

func NewAccionReforzar(jugador string, territorioReforzado NumRegion, tropasRefuerzo int) AccionReforzar {
	return AccionReforzar{
		IDAccion:            int(IDAccionReforzar),
		Jugador:             jugador,
		TerritorioReforzado: territorioReforzado,
		TropasRefuerzo:      tropasRefuerzo}
}

type AccionAtaque struct {
	IDAccion               int
	Origen                 NumRegion
	Destino                NumRegion
	TropasPerdidasAtacante NumRegion
	TropasPerdidasDefensor NumRegion
	NumDadosAtaque         int
	JugadorAtacante        string
	JugadorDefensor        string
}

func NewAccionAtaque(origen NumRegion, destino NumRegion, tropasPerdidasAtacante NumRegion, tropasPerdidasDefensor NumRegion, numDadosAtaque int, jugadorAtacante string, jugadorDefensor string) AccionAtaque {
	return AccionAtaque{IDAccion: int(IDAccionAtaque),
		Origen:                 origen,
		Destino:                destino,
		TropasPerdidasAtacante: tropasPerdidasAtacante,
		TropasPerdidasDefensor: tropasPerdidasDefensor,
		NumDadosAtaque:         numDadosAtaque,
		JugadorAtacante:        jugadorAtacante,
		JugadorDefensor:        jugadorDefensor}
}

type AccionOcupar struct {
	IDAccion        int
	Origen          NumRegion
	Destino         NumRegion
	TropasOrigen    NumRegion
	TropasDestino   NumRegion
	JugadorOcupante string
	JugadorOcupado  string
}

func NewAccionOcupar(origen NumRegion, destino NumRegion, tropasOrigen NumRegion, tropasDestino NumRegion, jugadorOcupante string, jugadorOcupado string) AccionOcupar {
	return AccionOcupar{
		IDAccion:        int(IDAccionOcupar),
		Origen:          origen,
		Destino:         destino,
		TropasOrigen:    tropasOrigen,
		TropasDestino:   tropasDestino,
		JugadorOcupante: jugadorOcupante,
		JugadorOcupado:  jugadorOcupado}
}

type AccionFortificar struct {
	IDAccion      int
	Origen        NumRegion
	Destino       NumRegion
	TropasOrigen  int
	TropasDestino int
	Jugador       string
}

func NewAccionFortificar(origen, destino NumRegion, tropasOrigen, tropasDestino int, jugador string) AccionFortificar {
	return AccionFortificar{
		IDAccion:      int(IDAccionFortificar),
		Origen:        origen,
		Destino:       destino,
		TropasOrigen:  tropasOrigen,
		TropasDestino: tropasDestino,
		Jugador:       jugador}
}

type AccionObtenerCarta struct {
	IDAccion int
	Carta    Carta
	Jugador  string
}

func NewAccionObtenerCarta(carta Carta, jugador string) AccionObtenerCarta {
	return AccionObtenerCarta{
		IDAccion: int(IDAccionObtenerCarta),
		Carta:    carta,
		Jugador:  jugador,
	}
}
