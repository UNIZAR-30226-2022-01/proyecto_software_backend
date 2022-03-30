// Package globales contiene variables globales a ser utilizadas por todos los módulos, instanciadas
// desde el paquete principal
package globales

import (
	"backend/vo"
	"database/sql" // Funciones de sql
	"sync"
)

const (
	DIRECCION_DB       = "DIRECCION_DB"
	DIRECCION_DB_TESTS = "DIRECCION_DB_TESTS"
	PUERTO_WEB         = "PUERTO_WEB"
	PUERTO_API         = "PUERTO_API"
	USUARIO_DB         = "USUARIO_DB"
	PASSWORD_DB        = "PASSWORD_DB"
	CARPETA_FRONTEND   = "web"
)

var Db *sql.DB // Base de datos thread safe, a compartir entre los módulos

var CachePartidas *AlmacenPartidas

type AlmacenPartidas struct {
	Mtx sync.RWMutex // Mutex 1 Escritor - N lectores

	Partidas map[int]vo.Partida

	CanalSerializacion chan vo.Partida
	CanalParada        chan struct{}
}

func IniciarAlmacenPartidas() *AlmacenPartidas {
	var ap AlmacenPartidas
	ap.Partidas = make(map[int]vo.Partida)
	ap.CanalSerializacion = make(chan vo.Partida, 50) // Estimación de partidas posibles a la vez
	ap.CanalParada = make(chan struct{})

	return &ap
}

// ObtenerPartida devuelve una copia de la partida con ID dado, y si existe o no
func (ap *AlmacenPartidas) ObtenerPartida(idp int) (partida vo.Partida, existe bool) {
	ap.Mtx.RLock()
	defer ap.Mtx.RUnlock()

	partida, existe = ap.Partidas[idp]

	return partida, existe
}

// AlmacenarPartida almacena o sobreescribe una partida en el almacén
func (ap *AlmacenPartidas) AlmacenarPartida(partida vo.Partida) {
	ap.Mtx.Lock()
	defer ap.Mtx.Unlock()

	ap.Partidas[partida.IdPartida] = partida
}

// EliminarPartida elimina una partida del almacén
func (ap *AlmacenPartidas) EliminarPartida(partida vo.Partida) {
	ap.Mtx.Lock()
	defer ap.Mtx.Unlock()
	delete(ap.Partidas, partida.IdPartida)
}

func (ap *AlmacenPartidas) PararAlmacenPartidas() {
	ap.CanalParada <- struct{}{}
}
