// Package globales contiene variables globales a ser utilizadas por todos los módulos, instanciadas
// desde el paquete principal
package globales

import (
	"database/sql" // Funciones de sql
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"log"
	"sync"
	"time"
)

const (
	DIRECCION_DB                   = "DIRECCION_DB"
	DIRECCION_DB_TESTS             = "DIRECCION_DB_TESTS"
	PUERTO_WEB                     = "PUERTO_WEB"
	PUERTO_API                     = "PUERTO_API"
	USUARIO_DB                     = "USUARIO_DB"
	PASSWORD_DB                    = "PASSWORD_DB"
	CARPETA_FRONTEND               = "web"
	INTERVALO_HORAS_LIMPIEZA_CACHE = 6
)

var Db *sql.DB // Base de datos thread safe, a compartir entre los módulos

var CachePartidas *AlmacenPartidas

var CanalEliminacionPartidasDB chan int        // Canal de eliminación de partidas con usuarios inactivos de la base de datos
var CanalExpulsionUsuariosDB chan string       // Canal de desvinculación de usuarios inactivos de sus partidas
var CanalParadaBorradoPartidasDB chan struct{} // Canal de parada de la Goroutine de atención a borrado de partidas y usuarios

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

	// Inicia una Goroutine de limpieza periódica de partidas a borrar sin intervención (partidas
	// con todos los jugadores expulsados)
	go func() {
		for {
			//time.Sleep(time.Second * 2) // Para pruebas
			time.Sleep(INTERVALO_HORAS_LIMPIEZA_CACHE * time.Hour)
			ap.limpiarCache()
		}
	}()

	return &ap
}

func IniciarCanalesEliminacionPartidasDB() {
	CanalEliminacionPartidasDB = make(chan int, 25) // Hasta 25 partidas siendo borradas concurrentemente, no se espera un tráfico cercano a este
	CanalParadaBorradoPartidasDB = make(chan struct{})
	CanalExpulsionUsuariosDB = make(chan string)
}

func (ap *AlmacenPartidas) limpiarCache() {
	ap.Mtx.Lock()
	defer ap.Mtx.Unlock()

	log.Println("Iniciando limpieza de cache...")
	for i, p := range ap.Partidas {
		//if p.Estado.UltimaAccion.Add(time.Second * 5).After(time.Now()) { // Para pruebas
		if p.Estado.UltimaAccion.Add(time.Hour * logica_juego.HORAS_EXPULSION_INACTIVIDAD).After(time.Now()) {
			log.Println("Eliminando a usuario", p.Estado.Jugadores[p.Estado.TurnoJugador], "de partida", i, "por inactividad...")
			CanalExpulsionUsuariosDB <- p.Estado.Jugadores[p.Estado.TurnoJugador]
			p.Estado.ExpulsarJugador()

			ap.Partidas[p.IdPartida] = p
		}

		if p.Estado.TerminadaPorExpulsiones() {
			log.Println("Borrando partida", p, "por estar todos los jugadores expulsados...")
			// La borra de la base de datos
			CanalEliminacionPartidasDB <- i
			// Y de la cache
			delete(ap.Partidas, p.IdPartida)
		}
	}
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

// EliminarPartida elimina una partida del almacén y se encarga del cierre correcto de sus goroutines asociadas
func (ap *AlmacenPartidas) EliminarPartida(partida vo.Partida) {
	ap.Mtx.Lock()
	defer ap.Mtx.Unlock()

	delete(ap.Partidas, partida.IdPartida)
}

func (ap *AlmacenPartidas) PararAlmacenPartidas() {
	ap.CanalParada <- struct{}{}
}
