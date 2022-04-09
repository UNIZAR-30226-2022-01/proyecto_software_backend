package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

// Prueba de funcionamiento de los subgrafos de regiones
// a partir del estado inicial de una partida, comprobando que
// cada subregión está formada por los territorios controlados
// de cada jugador
func TestSubgrafoRegiones(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuarios...")

	usuarios := []string{"usuario1", "usuario2", "usuario3", "usuario4", "usuario5", "usuario6"}

	cookie := crearUsuario(usuarios[0], t)
	cookie2 := crearUsuario(usuarios[1], t)
	cookie3 := crearUsuario(usuarios[2], t)
	cookie4 := crearUsuario(usuarios[3], t)
	cookie5 := crearUsuario(usuarios[4], t)
	cookie6 := crearUsuario(usuarios[5], t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, true)
	unirseAPartida(cookie2, t, 1)
	unirseAPartida(cookie3, t, 1)
	unirseAPartida(cookie4, t, 1)
	unirseAPartida(cookie5, t, 1)
	unirseAPartida(cookie6, t, 1)

	partidaCache := comprobarPartidaEnCurso(t, 1)

	regiones := make([][]logica_juego.NumRegion, 6)
	regiones[0] = obtenerRegionesSubGrafo(partidaCache, usuarios[0])
	regiones[1] = obtenerRegionesSubGrafo(partidaCache, usuarios[1])
	regiones[2] = obtenerRegionesSubGrafo(partidaCache, usuarios[2])
	regiones[3] = obtenerRegionesSubGrafo(partidaCache, usuarios[3])
	regiones[4] = obtenerRegionesSubGrafo(partidaCache, usuarios[4])
	regiones[5] = obtenerRegionesSubGrafo(partidaCache, usuarios[5])

	// Para cada región, buscarla en las subregiones. Solo debería estar en un único subgrafo
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		encontrada := false
		for indiceJugador, regionesSubGrafo := range regiones {
			if buscarEnLista(regionesSubGrafo, i) {
				if !encontrada {
					encontrada = true

					if partidaCache.Estado.EstadoMapa[i].Ocupante != usuarios[indiceJugador] {
						t.Fatal("La región", i, "está en el subgrafo de", usuarios[indiceJugador], "pero está controlada por", partidaCache.Estado.EstadoMapa[i].Ocupante)
					} else {
						t.Log("Región", i, "encontrada en subgrafo de", usuarios[indiceJugador], ", controlada por", partidaCache.Estado.EstadoMapa[i].Ocupante)
					}
				} else {
					t.Fatal("La región", i, "ya estaba en el subgrafo de otro jugador")
				}
			}
		}

		if !encontrada {
			t.Fatal("La región", i, "no se ha encontrada en el subgrafo de ningún jugador")
		}
	}
}
