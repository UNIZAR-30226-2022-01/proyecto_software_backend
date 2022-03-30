package logica_juego

import (
	"math/rand"
	"time"
)

type Carta struct {
	Tipo      TipoTropa
	Region    NumRegion
	EsComodin bool
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
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	rand.Shuffle(len(cartas), func(i, j int) { cartas[i], cartas[j] = cartas[j], cartas[i] })

	return cartas
}
