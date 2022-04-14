package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"net/http"
	"testing"
)

// Ataca hasta ocupar una región dada por el jugador con la cookie indicada, desde la región de origen dada
func atacarYOcupar(t *testing.T, partidaCache vo.Partida, cookies []*http.Cookie, regionOrigen logica_juego.NumRegion, regionDestino logica_juego.NumRegion) {
	saltarTurnos(t, partidaCache, "usuario1")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	pasarAFaseAtaque(partidaCache)
	partidaCache = comprobarPartidaEnCurso(t, 1)

	t.Log("turno:", partidaCache.Estado.Jugadores[partidaCache.Estado.TurnoJugador])
	t.Log("fase:", partidaCache.Estado.Fase)

	for partidaCache.Estado.EstadoMapa[regionDestino].Ocupante != "usuario1" {
		err := atacar(regionOrigen, regionDestino, 1, cookies[0], t)
		if err != nil {
			t.Fatal(err)
		}

		partidaCache = comprobarPartidaEnCurso(t, 1)

		if partidaCache.Estado.HayTerritorioDesocupado {
			err = ocupar(regionDestino, 1, cookies[0], t)
			if err != nil {
				t.Fatal(err)
			}
		}

		partidaCache = comprobarPartidaEnCurso(t, 1)
	}
}
