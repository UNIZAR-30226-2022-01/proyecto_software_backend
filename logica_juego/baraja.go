package logica_juego

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type Carta struct {
	IdCarta   int
	Tipo      TipoTropa
	Region    NumRegion
	EsComodin bool
}

// crearBaraja
func crearBaraja() (cartas []Carta) {
	idCarta := 0
	for i := Eastern_australia; i <= Alberta; i++ {
		var carta Carta
		var tipo TipoTropa
		if idCarta < 18 {
			tipo = Infanteria
		} else if idCarta < 36 {
			tipo = Caballeria
		} else {
			tipo = Artilleria
		}

		carta = Carta{
			IdCarta:   idCarta,
			Tipo:      tipo,
			Region:    i,
			EsComodin: false,
		}

		idCarta = idCarta + 1

		cartas = append(cartas, carta)
	}

	cartas = append(cartas, Carta{IdCarta: idCarta, EsComodin: true})
	cartas = append(cartas, Carta{IdCarta: idCarta + 1, EsComodin: true})

	barajarCartas(cartas)
	return cartas
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

	e.Acciones = append(e.Acciones, NewAccionCambioCartas(numTropas, hayBonificacion, regionBonificacion, numeroCartasInicial >= 5))
	return nil
}

// ConsultarCartas devuelve un slice que contiene las cartas que posee el usuario "jugador"
func (e *EstadoPartida) ConsultarCartas(jugador string) []Carta {
	return e.EstadosJugadores[jugador].Cartas
}

// retirarPrimeraCarta devuelve la primera carta del conjunto "cartas", o un error en caso de que no haya ninguna
func retirarPrimeraCarta(cartas []Carta) (carta Carta, cartasRes []Carta, err error) {
	if len(cartas) > 0 {
		carta = cartas[0]
		cartasRes = cartas[1:]

		return carta, cartasRes, nil
	}

	return Carta{}, []Carta{}, errors.New("El conjunto de cartas está vacío")
}

// RetirarCartaPorID retira la carta identificada por "id" del conjunto "cartas" en caso de que exista
// En cualquier otro caso devuelve un error. Usada para tests.
func RetirarCartaPorID(id int, cartas []Carta) (carta Carta, cartasRes []Carta, err error) {
	for i, c := range cartas {
		if c.IdCarta == id {
			cartasRes = append(cartas[0:i], cartas[i+1:]...)
			return c, cartasRes, nil
		}
	}

	return Carta{}, []Carta{}, errors.New("La carta con id " + strconv.Itoa(id) + " no existe en el conjunto")
}

// existeCarta devuelve true si existe la carta con identificador "id" en el conjunto "cartas", false si no
func existeCarta(id int, cartas []Carta) bool {
	for _, c := range cartas {
		if c.IdCarta == id {
			return true
		}
	}

	return false
}

// barajarCartas baraja de forma aleatoria un conjunto de cartas
func barajarCartas(cartas []Carta) {
	// Se baraja aleatoriamente
	rand.Seed(time.Now().UnixNano())
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	rand.Shuffle(len(cartas), func(i, j int) { cartas[i], cartas[j] = cartas[j], cartas[i] })
}

// esCambioValido devuelve true si se puede realizar un cambio correctamente con las cartas introducidas
// Será necesario que haya 3 cartas del mismo tipo o 2 cartas del mismo tipo además de un comodín o
// 3 cartas de cada tipo
func esCambioValido(cartas []Carta) bool {
	if len(cartas) != 3 {
		return false
	}

	comodines := 0
	infanteria := 0
	caballeria := 0
	artilleria := 0

	for _, c := range cartas {
		if c.EsComodin {
			comodines++
		} else if c.Tipo == Infanteria {
			infanteria++
		} else if c.Tipo == Caballeria {
			caballeria++
		} else {
			artilleria++
		}
	}

	return (comodines == 1 && (artilleria == 2 || caballeria == 2 || infanteria == 2)) ||
		artilleria == 3 || caballeria == 3 || infanteria == 3 ||
		(artilleria == 1 && caballeria == 1 && infanteria == 1)
}

func obtenerRegionesCartas(cartas []Carta) (regiones []NumRegion) {
	for _, c := range cartas {
		regiones = append(regiones, c.Region)
	}
	return regiones
}
