package logica_juego

import "encoding/gob"

const (
	IDAccionRecibirRegion Fase = iota
	IDAccionCambioFase
	IDAccionInicioTurno
	IDAccionCambioCartas
	IDAccionReforzar
	IDAccionAtaque
	IDAccionOcupar
	IDAccionFortificar
	IDAccionObtenerCarta
)

// AccionRecibirRegion corresponde a la asignación automática de un territorio a
// un jugador dado durante el inicio de la partida.
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 0,
// 		"Region": 1,
// 		"TropasRestantes": 4,
// 		"TerritoriosRestantes": 8,
// 		"Jugador": "usuario1"
//    }
type AccionRecibirRegion struct {
	IDAccion             int       // 0
	Region               NumRegion // Región asignada
	TropasRestantes      int       // Tropas que tiene el jugador una vez asignado el territorio
	TerritoriosRestantes int       // Territorios restantes en el mapa sin asignar
	Jugador              string    // Nombre de jugador receptor del territorio
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

// AccionCambioFase corresponde a un cambio de fase dentro del turno del jugador dado.
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 1,
//      "Fase": 2,
//      "Jugador": "usuario1"
//    }
type AccionCambioFase struct {
	IDAccion int    // 1
	Fase     Fase   // {0: Inicio (no usada), 1: Refuerzo , 2: Ataque , 3: Fortificar}
	Jugador  string // Jugador del turno
}

func NewAccionCambioFase(fase Fase, jugador string) AccionCambioFase {
	return AccionCambioFase{
		IDAccion: int(IDAccionCambioFase),
		Fase:     fase,
		Jugador:  jugador}
}

// AccionInicioTurno corresponde a un cambio de turno al usuario dado. El resto de usuarios no tendrán éxito en peticiones
// que no sean de solicitud de estado durante su turno.
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 2,
// 		"Jugador": "usuario1",
// 		"TropasObtenidas": 2,
// 		"RazonNumeroTerritorios": 12,
// 		"RazonContinentesOcupados": 1
//    }
type AccionInicioTurno struct {
	IDAccion                 int    // 2
	Jugador                  string // Jugador del nuevo turno
	TropasObtenidas          int    // Tropas obtenidas durante el nuevo turno
	RazonNumeroTerritorios   int    // Número de territorios debido a los cuales ha recibido dicho número de tropas
	RazonContinentesOcupados int    // Número de continentes debido a los cuales ha recibido dicho número de tropas
}

func NewAccionInicioTurno(jugador string, tropasObtenidas int, razonNumeroTerritorios int, razonContinentesOcupados int) AccionInicioTurno {
	return AccionInicioTurno{
		IDAccion:                 int(IDAccionInicioTurno),
		Jugador:                  jugador,
		TropasObtenidas:          tropasObtenidas,
		RazonNumeroTerritorios:   razonNumeroTerritorios,
		RazonContinentesOcupados: razonContinentesOcupados}
}

// AccionCambioCartas corresponde a un cambio de turno al usuario dado. El resto de usuarios no tendrán éxito en peticiones
// que no sean de solicitud de estado durante su turno.
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 3,
// 		"NumConjuntosCambiados": 1,
// 		"NumTropasObtenidas": 2,
// 		"BonificacionObtenida": true,
// 		"RegionQueOtorgaBonificacion": 2,
// 		"ObligadoAHacerCambios": false
//    }
type AccionCambioCartas struct {
	IDAccion                    int       // 3
	NumConjuntosCambiados       int       // Número de cartas cambiadas (múltiplos de 3)
	NumTropasObtenidas          int       // Tropas obtenidas por el cambio
	BonificacionObtenida        bool      // Flag de si se ha recibido una bonificación de territorio de una de las cartas usadas
	RegionQueOtorgaBonificacion NumRegion // ID de región que ha otorgado la bonificación, si se ha obtenido
	ObligadoAHacerCambios       bool      // Flag de si el usuario ha sido obligado a hacer el cambio, por tener más de 5 cartas
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

// AccionReforzar corresponde a un refuerzo de una región por un jugador
//
// Ejemplo en JSON:
//    {
//		"IDAccion": 4,
//		"Jugador": "usuario1",
//		"TerritorioReforzado": 1,
//		"TropasRefuerzo": 20
//    }
type AccionReforzar struct {
	IDAccion            int       // 4
	Jugador             string    // Jugador que ha reforzado el territorio
	TerritorioReforzado NumRegion // ID de región que ha sido reforzada
	TropasRefuerzo      int       // Número de tropas de refuerzo asignadas a la región
}

func NewAccionReforzar(jugador string, territorioReforzado NumRegion, tropasRefuerzo int) AccionReforzar {
	return AccionReforzar{
		IDAccion:            int(IDAccionReforzar),
		Jugador:             jugador,
		TerritorioReforzado: territorioReforzado,
		TropasRefuerzo:      tropasRefuerzo}
}

// AccionAtaque corresponde al ataque de una región por parte de un usuario dado
//
// Ejemplo en JSON:
//    {
//  	"IDAccion": 5,
//  	"Origen": 2,
//  	"Destino": 3,
//  	"TropasPerdidasAtacante": 15,
//  	"TropasPerdidasDefensor": 5,
//  	"NumDadosAtaque": 3,
//  	"JugadorAtacante": "usuario1",
//  	"JugadorDefensor": "usuario2"
//    }
type AccionAtaque struct {
	IDAccion               int       // 5
	Origen                 NumRegion // ID de región de la cual se origina el ataque (y usan sus tropas)
	Destino                NumRegion // ID de región atacada
	TropasPerdidasAtacante NumRegion // Tropas perdidas por el atacante
	TropasPerdidasDefensor NumRegion // Tropas perdidas por el defensor
	NumDadosAtaque         int       // Número de dados lanzados por el atacante
	JugadorAtacante        string    // Nombre del atacante
	JugadorDefensor        string    // Nombre del defensor
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

// AccionOcupar corresponde a la ocupación de una región por un jugador tras un ataque con éxito.
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 6,
// 		"Origen": 2,
// 		"Destino": 3,
// 		"TropasOrigen": 10,
// 		"TropasDestino": 5,
// 		"JugadorOcupante": "usuario1",
// 		"JugadorOcupado": "usuario2"
//    }
type AccionOcupar struct {
	IDAccion        int       // 6
	Origen          NumRegion // ID de región desde la cual se originó el ataque (y usaron sus tropas)
	Destino         NumRegion // ID de región ocupada
	TropasOrigen    NumRegion // Número de tropas que han quedado en la región desde la cual se originó el ataque
	TropasDestino   NumRegion // Número de tropas asignadas a la región ocupada
	JugadorOcupante string    // Nombre del atacante
	JugadorOcupado  string    // Nombre del defensor que ha perdido el territorio
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

// AccionFortificar corresponde a la fortificación de un territorio de un jugador
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 7,
// 		"Origen": 7,
// 		"Destino": 9,
// 		"TropasOrigen": 10,
// 		"TropasDestino": 8,
// 		"Jugador": "usuario1"
//    }
type AccionFortificar struct {
	IDAccion      int       // 7
	Origen        NumRegion // ID de región desde la cual se han movido tropas
	Destino       NumRegion // ID de región que ha recibido las tropas
	TropasOrigen  int       // Número de tropas que han quedado en la región desde la cual se originó el movimiento
	TropasDestino int       // Número de tropas que hay en la región fortificada tras la acción
	Jugador       string    // Jugador que ha fortificado
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

// AccionObtenerCarta corresponde a la recepción de una carta por parte de un jugador
//
// Ejemplo en JSON:
//    {
// 		"IDAccion": 8,
// 		"Carta": {
//          		"Tipo": 0,					// {0: Infanteria, 1: Caballeria, 2: Artilleria}
//          		"Region": 29,				// ID de región de la carta
//          		"EsComodin": false			// Flag de si la carta es un comodín (tiene cualquier tipo de tropa y no tiene región)
//          	 },
//  	"Jugador": "usuario1"
//    }
type AccionObtenerCarta struct {
	IDAccion int    // 8
	Carta    Carta  // Carta recibida
	Jugador  string // Jugador receptor
}

func NewAccionObtenerCarta(carta Carta, jugador string) AccionObtenerCarta {
	return AccionObtenerCarta{
		IDAccion: int(IDAccionObtenerCarta),
		Carta:    carta,
		Jugador:  jugador,
	}
}

func RegistrarAcciones() {
	gob.Register(AccionRecibirRegion{})
	gob.Register(AccionCambioFase{})
	gob.Register(AccionInicioTurno{})
	gob.Register(AccionCambioCartas{})
	gob.Register(AccionReforzar{})
	gob.Register(AccionAtaque{})
	gob.Register(AccionOcupar{})
	gob.Register(AccionFortificar{})
	gob.Register(AccionObtenerCarta{})
	gob.Register(struct{}{}) // Placeholder de acciones no implementadas
}
