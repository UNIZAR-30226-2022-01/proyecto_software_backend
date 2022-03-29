// Package globales contiene variables globales a ser utilizadas por todos los módulos, instanciadas
// desde el paquete principal
package globales

import (
	"backend/vo"
	"database/sql" // Funciones de sql
	"sync"

	grafos "gonum.org/v1/gonum/graph/simple"
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

var GrafoMapa *grafos.UndirectedGraph

func InicializarGrafoMapa() {
	GrafoMapa = grafos.NewUndirectedGraph()

	for i := vo.Eastern_australia; i < vo.Alberta; i++ {
		node := grafos.Node(i)
		GrafoMapa.AddNode(node)
	}
	// América del norte
	añadirArista(vo.Alaska, vo.Northwest_territory)
	añadirArista(vo.Alaska, vo.Alberta)
	añadirArista(vo.Alberta, vo.Northwest_territory)
	añadirArista(vo.Alberta, vo.Ontario)
	añadirArista(vo.Northwest_territory, vo.Ontario)
	añadirArista(vo.Northwest_territory, vo.Greenland)
	añadirArista(vo.Greenland, vo.Ontario)
	añadirArista(vo.Quebec, vo.Greenland)
	añadirArista(vo.Quebec, vo.Ontario)
	añadirArista(vo.Quebec, vo.Eastern_united_states)
	añadirArista(vo.Alberta, vo.Western_united_states)
	añadirArista(vo.Ontario, vo.Eastern_united_states)
	añadirArista(vo.Ontario, vo.Western_united_states)
	añadirArista(vo.Western_united_states, vo.Eastern_united_states)
	añadirArista(vo.Central_america, vo.Western_united_states)
	añadirArista(vo.Central_america, vo.Eastern_united_states)

	// América del sur
	añadirArista(vo.Venezuela, vo.Brazil)
	añadirArista(vo.Venezuela, vo.Peru)
	añadirArista(vo.Peru, vo.Brazil)
	añadirArista(vo.Peru, vo.Argentina)
	añadirArista(vo.Argentina, vo.Brazil)

	// Europa
	añadirArista(vo.Iceland, vo.Scandinavia)
	añadirArista(vo.Iceland, vo.Great_britain)

	añadirArista(vo.Scandinavia, vo.Great_britain)
	añadirArista(vo.Scandinavia, vo.Northern_europe)
	añadirArista(vo.Scandinavia, vo.Ukraine)

	añadirArista(vo.Great_britain, vo.Northern_europe)
	añadirArista(vo.Great_britain, vo.Western_europe)

	añadirArista(vo.Northern_europe, vo.Ukraine)
	añadirArista(vo.Northern_europe, vo.Western_europe)
	añadirArista(vo.Northern_europe, vo.Southern_europe)

	añadirArista(vo.Ukraine, vo.Southern_europe)

	añadirArista(vo.Western_europe, vo.Southern_europe)

	// África
	añadirArista(vo.North_africa, vo.Egypt)
	añadirArista(vo.North_africa, vo.East_africa)
	añadirArista(vo.North_africa, vo.Congo)

	añadirArista(vo.Egypt, vo.East_africa)

	añadirArista(vo.Congo, vo.East_africa)
	añadirArista(vo.Congo, vo.South_africa)

	añadirArista(vo.East_africa, vo.South_africa)
	añadirArista(vo.East_africa, vo.Madagascar)
	añadirArista(vo.South_africa, vo.Madagascar)

	// Asia
	añadirArista(vo.Yakursk, vo.Siberia)
	añadirArista(vo.Yakursk, vo.Irkutsk)
	añadirArista(vo.Yakursk, vo.Kamchatka)
	añadirArista(vo.Ural, vo.Siberia)
	añadirArista(vo.Ural, vo.China)
	añadirArista(vo.Ural, vo.Afghanistan)
	añadirArista(vo.Siberia, vo.Irkutsk)
	añadirArista(vo.Siberia, vo.Mongolia)
	añadirArista(vo.Siberia, vo.China)
	añadirArista(vo.Irkutsk, vo.Kamchatka)
	añadirArista(vo.Irkutsk, vo.Mongolia)
	añadirArista(vo.Kamchatka, vo.Mongolia)
	añadirArista(vo.Kamchatka, vo.Japan)
	añadirArista(vo.Afghanistan, vo.China)
	añadirArista(vo.Afghanistan, vo.India)
	añadirArista(vo.Afghanistan, vo.Middle_east)
	añadirArista(vo.China, vo.Mongolia)
	añadirArista(vo.China, vo.Siam)
	añadirArista(vo.China, vo.India)
	añadirArista(vo.Mongolia, vo.Japan)
	añadirArista(vo.Middle_east, vo.India)
	añadirArista(vo.India, vo.Siam)

	// Australia
	añadirArista(vo.Indonesia, vo.New_guinea)
	añadirArista(vo.Indonesia, vo.Western_australia)
	añadirArista(vo.New_guinea, vo.Eastern_australia)
	añadirArista(vo.New_guinea, vo.Western_australia)
	añadirArista(vo.Western_australia, vo.Eastern_australia)

	// Conexiones entre continentes
	añadirArista(vo.Alaska, vo.Kamchatka)
	añadirArista(vo.Greenland, vo.Iceland)

	añadirArista(vo.Brazil, vo.North_africa)
	añadirArista(vo.Venezuela, vo.Central_america)

	añadirArista(vo.Ukraine, vo.Ural)
	añadirArista(vo.Ukraine, vo.Afghanistan)
	añadirArista(vo.Ukraine, vo.Middle_east)

	añadirArista(vo.Southern_europe, vo.Middle_east)
	añadirArista(vo.Southern_europe, vo.North_africa)
	añadirArista(vo.Southern_europe, vo.Egypt)
	añadirArista(vo.Western_europe, vo.North_africa)
	añadirArista(vo.Egypt, vo.Middle_east)
	añadirArista(vo.East_africa, vo.Middle_east)
	añadirArista(vo.Indonesia, vo.Siam)
}

func añadirArista(region1 vo.NumRegion, region2 vo.NumRegion) {
	nodo1 := grafos.Node(region1)
	nodo2 := grafos.Node(region2)
	GrafoMapa.SetEdge(GrafoMapa.NewEdge(nodo1, nodo2))
}

func Conectadas(region1 vo.NumRegion, region2 vo.NumRegion) bool {
	nodo1 := grafos.Node(region1)
	nodo2 := grafos.Node(region2)

	return GrafoMapa.HasEdgeBetween(nodo1.ID(), nodo2.ID())
}

var CachePartidas *AlmacenPartidas

type AlmacenPartidas struct {
	Mtx sync.RWMutex // Mutex 1 Escritor - N lectores

	Partidas map[int]vo.Partida

	CanalSerializacion chan vo.Partida
	CanalParada        chan struct{}
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
