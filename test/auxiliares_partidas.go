package integracion

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func saltarTurnos(t *testing.T, partidaCache vo.Partida, usuario string) {
	t.Log("Turno actual:", partidaCache.Estado.ObtenerJugadorTurno())
	t.Log("Saltando turnos hasta " + usuario + "...")

	for partidaCache.Estado.ObtenerJugadorTurno() != usuario {

		partidaCache = comprobarPartidaEnCurso(t, 1)
		t.Log("Turno saltado:", partidaCache.Estado.ObtenerJugadorTurno())

		partidaCache.Estado.SiguienteJugador()
		globales.CachePartidas.AlmacenarPartida(partidaCache)
	}

	t.Log("Turno nuevo:", partidaCache.Estado.ObtenerJugadorTurno())
}

func comprobarPartidaNoEnCurso(t *testing.T, idp int) vo.Partida {
	partidaCache, existe := globales.CachePartidas.ObtenerPartida(idp)
	if !existe {
		t.Fatal("No hay partidas en la Cache, aunque debería haber.")
	} else if partidaCache.EnCurso {
		t.Fatal("La partida está en curso, aunque no debería.")
	}

	return partidaCache
}

func comprobarPartidaEnCurso(t *testing.T, idp int) vo.Partida {
	partidaCache, existe := globales.CachePartidas.ObtenerPartida(idp)
	if !existe {
		t.Fatal("No hay partidas en la Cache, aunque debería haber.")
	} else if !partidaCache.EnCurso {
		t.Fatal("La partida no está en curso, aunque debería.")
	}

	return partidaCache
}

func crearPartida(cookie *http.Cookie, t *testing.T, publica bool) {
	// O usar cookie jar de https://stackoverflow.com/questions/12756782/go-http-post-and-use-cookies
	client := &http.Client{}

	var campos url.Values
	if publica {
		campos = url.Values{
			"password":     {""},
			"maxJugadores": {"6"},
			"tipo":         {"Publica"}, // o "Privada"
		}
	} else {
		campos = url.Values{
			"password":     {"password"},
			"maxJugadores": {"6"},
			"tipo":         {"Privada"}, // o "Privada"
		}
	}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/crearPartida", strings.NewReader(campos.Encode()))
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de crear partida:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al crear una partida:", resp.StatusCode)
	}

	log.Println("partida creada")
}

// Idéntica a la anterior, pero con 3 jugadores
func crearPartidaReducida(cookie *http.Cookie, t *testing.T, publica bool) {
	client := &http.Client{}

	var campos url.Values
	if publica {
		campos = url.Values{
			"password":     {""},
			"maxJugadores": {"3"},
			"tipo":         {"Publica"}, // o "Privada"
		}
	} else {
		campos = url.Values{
			"password":     {"password"},
			"maxJugadores": {"3"},
			"tipo":         {"Privada"}, // o "Privada"
		}
	}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/crearPartida", strings.NewReader(campos.Encode()))
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de crear partida:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al crear una partida:", resp.StatusCode)
	}
}

func consultarEstadoLobby(cookie *http.Cookie, t *testing.T) (estado vo.EstadoLobby) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerEstadoLobby", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en GET de consultar lobby:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al consultar lobby:", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&estado)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al consultar el lobby:", err)
		}

		t.Log("Respuesta de consultar lobby:", estado)
		return estado
	}

	return estado
}

func abandonarLobby(cookie *http.Cookie, t *testing.T) {
	// O usar cookie jar de https://stackoverflow.com/questions/12756782/go-http-post-and-use-cookies
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/abandonarLobby", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de abandonar partida:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al abandonar una partida:", resp.StatusCode)
	}
}

func abandonarPartida(cookie *http.Cookie, t *testing.T) {
	// O usar cookie jar de https://stackoverflow.com/questions/12756782/go-http-post-and-use-cookies
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/abandonarPartida", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de abandonar partida:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al abandonar una partida:", resp.StatusCode)
	}
}

func unirseAPartida(cookie *http.Cookie, t *testing.T, id int) {
	client := &http.Client{}

	campos := url.Values{
		"idPartida": {strconv.Itoa(id)},
		"password":  {"password"},
	}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/unirseAPartida", strings.NewReader(campos.Encode()))
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de unirse a partida:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al unirse a una partida:", resp.StatusCode)
	}
}

func jugandoEnPartida(cookie *http.Cookie, t *testing.T) (esta bool) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/jugandoEnPartida", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en GET de jugandoEnPartida:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 de jugandoEnPartida:", resp.StatusCode)
	} else {
		err = json.NewDecoder(resp.Body).Decode(&esta)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta de jugandoEnPartida:", err)
		}

		t.Log("Respuesta de jugandoEnPartida:", esta)
		return esta
	}

	return esta
}

func obtenerPartidas(cookie *http.Cookie, t *testing.T) []vo.ElementoListaPartidas {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerPartidas", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en GET de obtener partidas:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al obtener partidas:", resp.StatusCode)
	} else {
		var partidas []vo.ElementoListaPartidas
		err = json.NewDecoder(resp.Body).Decode(&partidas)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al obtener partidas:", err)
		}

		t.Log("Respuesta de obtenerPartidas:", partidas)
		return partidas
	}

	return nil
}

func obtenerPartidaDB(t *testing.T, idP int) vo.Partida {
	partidaDB, err := dao.ObtenerPartida(globales.Db, idP)

	if err != nil {
		t.Fatal("Error al obtener partida de DB:", idP)
	}
	return partidaDB
}

// Fuerza la conquista de un territorio por un jugador dados en una partida de cache,
// almacenando el estado en el cache de vuelta
func conquistar(t *testing.T, partidaCache vo.Partida, territorio int, usuario string) vo.Partida {
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(territorio)].Ocupante = usuario

	globales.CachePartidas.AlmacenarPartida(partidaCache)
	return partidaCache
}

// Fuerza la asignación de tropas a un territorio por un jugador dados en una partida de cache,
// almacenando el estado en el cache de vuelta
func darTropas(t *testing.T, partidaCache vo.Partida, numTropas int, usuario string) vo.Partida {
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadosJugadores[usuario].Tropas += numTropas
	globales.CachePartidas.AlmacenarPartida(partidaCache)
	return partidaCache
}
