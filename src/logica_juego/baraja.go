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

// retirarPrimeraCarta devuelve la primera carta del conjunto "cartas", o un error en caso de que no haya ninguna
func retirarPrimeraCarta(cartas []Carta) (carta Carta, err error) {
	if len(cartas) > 0 {
		carta = cartas[0]
		cartas = cartas[1:]
		return carta, nil
	}

	return Carta{}, errors.New("El conjunto de cartas está vacío")
}

// retirarCartaPorID retira la carta identificada por "id" del conjunto "cartas" en caso de que exista
// En cualquier otro caso devuelve un error
func retirarCartaPorID(id int, cartas []Carta) (carta Carta, err error) {
	for i, c := range cartas {
		if c.IdCarta == id {
			cartas = append(cartas[0:i], cartas[i+1:]...)
			return c, nil
		}
	}

	return Carta{}, errors.New("La carta con id " + strconv.Itoa(id) + " no existe en el conjunto")
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
