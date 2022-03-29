package logica_juego

import (
	"math/rand"
	"sync"
	"time"
)

var onlyOnce sync.Once

func LanzarDados() int {
	var dados = []int{1, 2, 3, 4, 5, 6}

	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	return dados[rand.Intn(len(dados))] // Devuelve una posición [0, 6)
}
