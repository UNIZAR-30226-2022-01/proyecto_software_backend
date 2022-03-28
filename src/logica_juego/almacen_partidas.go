package logica_juego

import (
	//"backend/dao"

	"backend/vo"
	"sync"
)

type AlmacenPartidas struct {
	Mtx sync.RWMutex

	Partidas map[int]vo.Partida

	CanalSerializacion chan vo.Partida
	CanalParada        chan struct{}
}

func (ap *AlmacenPartidas) ObtenerPartida(idp int) (partida vo.Partida, existe bool) {
	ap.Mtx.RLock()
	defer ap.Mtx.RUnlock()

	partida, existe = ap.Partidas[idp]

	return partida, existe
}

func (ap *AlmacenPartidas) AlmacenarPartida(partida vo.Partida) {
	ap.Mtx.Lock()
	defer ap.Mtx.Unlock()

	ap.Partidas[partida.IdPartida] = partida
}

func (ap *AlmacenPartidas) EliminarPartida(partida vo.Partida) {
	ap.Mtx.Lock()
	defer ap.Mtx.Unlock()
	delete(ap.Partidas, partida.IdPartida)
}

func IniciarAlmacenPartidas() *AlmacenPartidas {
	var ap AlmacenPartidas
	ap.Partidas = make(map[int]vo.Partida)
	ap.CanalSerializacion = make(chan vo.Partida, 50) // Estimación de partidas posibles a la vez
	ap.CanalParada = make(chan struct{})

	return &ap
}

func (ap *AlmacenPartidas) PararAlmacenPartidas() {
	ap.CanalParada <- struct{}{}
}
