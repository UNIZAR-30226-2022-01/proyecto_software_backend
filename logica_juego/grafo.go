package logica_juego

import (
	"gonum.org/v1/gonum/graph/simple"
)

var GrafoMapa *simple.UndirectedGraph

func InicializarGrafoMapa() {
	GrafoMapa = simple.NewUndirectedGraph()

	for i := Eastern_australia; i < Alberta; i++ {
		node := simple.Node(i)
		GrafoMapa.AddNode(node)
	}
	// América del norte
	añadirArista(Alaska, Northwest_territory)
	añadirArista(Alaska, Alberta)
	añadirArista(Alberta, Northwest_territory)
	añadirArista(Alberta, Ontario)
	añadirArista(Northwest_territory, Ontario)
	añadirArista(Northwest_territory, Greenland)
	añadirArista(Greenland, Ontario)
	añadirArista(Quebec, Greenland)
	añadirArista(Quebec, Ontario)
	añadirArista(Quebec, Eastern_united_states)
	añadirArista(Alberta, Western_united_states)
	añadirArista(Ontario, Eastern_united_states)
	añadirArista(Ontario, Western_united_states)
	añadirArista(Western_united_states, Eastern_united_states)
	añadirArista(Central_america, Western_united_states)
	añadirArista(Central_america, Eastern_united_states)

	// América del sur
	añadirArista(Venezuela, Brazil)
	añadirArista(Venezuela, Peru)
	añadirArista(Peru, Brazil)
	añadirArista(Peru, Argentina)
	añadirArista(Argentina, Brazil)

	// Europa
	añadirArista(Iceland, Scandinavia)
	añadirArista(Iceland, Great_britain)

	añadirArista(Scandinavia, Great_britain)
	añadirArista(Scandinavia, Northern_europe)
	añadirArista(Scandinavia, Ukraine)

	añadirArista(Great_britain, Northern_europe)
	añadirArista(Great_britain, Western_europe)

	añadirArista(Northern_europe, Ukraine)
	añadirArista(Northern_europe, Western_europe)
	añadirArista(Northern_europe, Southern_europe)

	añadirArista(Ukraine, Southern_europe)

	añadirArista(Western_europe, Southern_europe)

	// África
	añadirArista(North_africa, Egypt)
	añadirArista(North_africa, East_africa)
	añadirArista(North_africa, Congo)

	añadirArista(Egypt, East_africa)

	añadirArista(Congo, East_africa)
	añadirArista(Congo, South_africa)

	añadirArista(East_africa, South_africa)
	añadirArista(East_africa, Madagascar)
	añadirArista(South_africa, Madagascar)

	// Asia
	añadirArista(Yakursk, Siberia)
	añadirArista(Yakursk, Irkutsk)
	añadirArista(Yakursk, Kamchatka)
	añadirArista(Ural, Siberia)
	añadirArista(Ural, China)
	añadirArista(Ural, Afghanistan)
	añadirArista(Siberia, Irkutsk)
	añadirArista(Siberia, Mongolia)
	añadirArista(Siberia, China)
	añadirArista(Irkutsk, Kamchatka)
	añadirArista(Irkutsk, Mongolia)
	añadirArista(Kamchatka, Mongolia)
	añadirArista(Kamchatka, Japan)
	añadirArista(Afghanistan, China)
	añadirArista(Afghanistan, India)
	añadirArista(Afghanistan, Middle_east)
	añadirArista(China, Mongolia)
	añadirArista(China, Siam)
	añadirArista(China, India)
	añadirArista(Mongolia, Japan)
	añadirArista(Middle_east, India)
	añadirArista(India, Siam)

	// Australia
	añadirArista(Indonesia, New_guinea)
	añadirArista(Indonesia, Western_australia)
	añadirArista(New_guinea, Eastern_australia)
	añadirArista(New_guinea, Western_australia)
	añadirArista(Western_australia, Eastern_australia)

	// Conexiones entre continentes
	añadirArista(Alaska, Kamchatka)
	añadirArista(Greenland, Iceland)

	añadirArista(Brazil, North_africa)
	añadirArista(Venezuela, Central_america)

	añadirArista(Ukraine, Ural)
	añadirArista(Ukraine, Afghanistan)
	añadirArista(Ukraine, Middle_east)

	añadirArista(Southern_europe, Middle_east)
	añadirArista(Southern_europe, North_africa)
	añadirArista(Southern_europe, Egypt)
	añadirArista(Western_europe, North_africa)
	añadirArista(Egypt, Middle_east)
	añadirArista(East_africa, Middle_east)
	añadirArista(Indonesia, Siam)
}

func añadirArista(region1 NumRegion, region2 NumRegion) {
	nodo1 := simple.Node(region1)
	nodo2 := simple.Node(region2)
	GrafoMapa.SetEdge(GrafoMapa.NewEdge(nodo1, nodo2))
}

func Conectadas(region1 NumRegion, region2 NumRegion) bool {
	nodo1 := simple.Node(region1)
	nodo2 := simple.Node(region2)

	return GrafoMapa.HasEdgeBetween(nodo1.ID(), nodo2.ID())
}

// Adyacentes devuelve una lista de regiones aydacentes a "region" en el grafo
func Adyacentes(region NumRegion) (regiones []NumRegion) {
	for i := Eastern_australia; i < Alberta; i++ {
		if Conectadas(region, i) {
			regiones = append(regiones, i)
		}
	}
	return regiones
}
