package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
)

// Devuelve una lista de regiones contenidas en el subgrafo del usuario dado
func obtenerRegionesSubGrafo(partidaCache vo.Partida, usuario string) (regiones []logica_juego.NumRegion) {
	grafo := partidaCache.Estado.ObtenerSubgrafoRegiones(usuario)

	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		if grafo.Node(int64(i)) != nil {
			regiones = append(regiones, i)
		}
	}

	return regiones
}
