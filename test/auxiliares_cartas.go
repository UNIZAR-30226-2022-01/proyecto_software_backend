package integracion

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"net/http"
	"os"
	"testing"
)

func cambiarCartas(t *testing.T, estadoJugador *logica_juego.EstadoJugador, estadoPartida *logica_juego.EstadoPartida, id1, id2, id3, numCanje int) {
	tropasIniciales := estadoJugador.Tropas
	numCartasInicial := len(estadoJugador.Cartas)
	var tropasEsperadas int
	if numCanje < 6 {
		tropasEsperadas = 4 + (numCanje-1)*2
	} else {
		tropasEsperadas = 15 + (numCanje-6)*5
	}
	err := estadoPartida.CambiarCartas("Jugador1", id1, id2, id3)
	if err != nil {
		t.Fatal("Error al cambiar 3 cartas:", err)
	}

	if (estadoJugador.Tropas - tropasIniciales) != tropasEsperadas {
		t.Fatal("El jugador debería recibir", tropasEsperadas, "tropas por el canje, pero recibe", estadoJugador.Tropas-tropasIniciales)
	}

	if (numCartasInicial - len(estadoJugador.Cartas)) != 3 {
		t.Fatal("Se deberían haber cambiado 3 cartas, pero se han cambiado:", numCartasInicial-len(estadoJugador.Cartas))
	}

	t.Log("Canje número:", numCanje, ";Se han recibido", estadoJugador.Tropas-tropasIniciales, "tropas a cambio de", numCartasInicial-len(estadoJugador.Cartas), "cartas")
}

func consultarCartas(cookie *http.Cookie, t *testing.T) (cartas []logica_juego.Carta) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/consultarCartas", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en GET de consultar cartas:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al consultar cartas:", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&cartas)
	}
	return cartas
}

// invarianteNumeroDeCartas comprueba que la suma de cartas de la baraja, descartes y mano del jugador es 44
func invarianteNumeroDeCartas(eP logica_juego.EstadoPartida, eJ logica_juego.EstadoJugador, t *testing.T) {
	t.Log("Cartas en la mano del jugador:", len(eJ.Cartas), ",cartas en la baraja:", len(eP.Cartas),
		",cartas descartadas:", len(eP.Descartes))
	cartas := len(eJ.Cartas) + len(eP.Cartas) + len(eP.Descartes)
	if cartas != 44 {
		t.Fatal("Hay un total de", cartas, "cartas, no se cumple el invariante")
	}
}

func cambioDeFaseConDemasiadasCartas(t *testing.T, partidaCache vo.Partida, err error, cookie *http.Cookie, usuario string) (vo.Partida, error) {
	partidaCache = comprobarPartidaEnCurso(t, 1)
	baraja := partidaCache.Estado.Cartas

	// Le damos al jugador las 6 primeras cartas (todas de infantería)
	for numCartas := 0; numCartas <= 5; numCartas++ {
		for _, carta := range baraja {
			if carta.IdCarta == numCartas {
				partidaCache.Estado.EstadosJugadores[usuario].Cartas = append(partidaCache.Estado.EstadosJugadores[usuario].Cartas, carta)
				break
			}
		}
	}

	globales.CachePartidas.AlmacenarPartida(partidaCache)
	err = saltarFase(cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al intentar saltar la fase")
	}
	t.Log("OK: No se ha podido saltar la fase, error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadosJugadores[usuario].Cartas = nil
	globales.CachePartidas.AlmacenarPartida(partidaCache)
	return partidaCache, err
}

func robarBarajaCompleta(e *logica_juego.EstadoPartida, t *testing.T) {
	e.HaConquistado = true
	e.HaRecibidoCarta = false
	err := e.RecibirCarta("Jugador1")
	for err == nil {
		e.HaConquistado = true
		e.HaRecibidoCarta = false
		err = e.RecibirCarta("Jugador1")
	}
	t.Log("Ya no quedan cartas, o error:", err)
}
