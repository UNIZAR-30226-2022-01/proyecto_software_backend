package logica_juego

import (
	"math/rand"
	"sync"
	"time"
)

var onlyOnce sync.Once

// LanzarDados devuelve un número entre [0-5] (utilizado para indexar,
// que se debe aumentar en una unidad para
func LanzarDados() int {
	var dados = []int{0, 1, 2, 3, 4, 5}

	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	return dados[rand.Intn(len(dados))] // Devuelve una posición [0, 6)
}
