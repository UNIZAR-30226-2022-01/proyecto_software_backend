package integracion

import (
	"backend/globales"
	"backend/logica_juego"
	"backend/middleware"
	"backend/vo"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func purgarDB() {
	_, err := globales.Db.Exec(`DELETE FROM "backend"."Partida"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."EsAmigo"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."Participa"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."TieneItems"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."Usuario"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`ALTER SEQUENCE backend."Partida_id_seq" RESTART`)
	if err != nil {
		log.Fatal(err)
	}
}

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

func reforzarTerritorioConFallo(t *testing.T, cookie *http.Cookie, numRegion int, numTropas int) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/reforzarTerritorio/"+strconv.Itoa(numRegion)+"/"+strconv.Itoa(numTropas), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace e

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en POST de reforzar territorio:", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatal("Obtenido código de error OK al forzar error en reforzar territorio:", resp.StatusCode)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		t.Log("Recibido error correctamente: " + bodyString)
	}

}

func reforzarTerritorio(t *testing.T, cookie *http.Cookie, numRegion int, numTropas int) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/reforzarTerritorio/"+strconv.Itoa(numRegion)+"/"+strconv.Itoa(numTropas), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace e

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en POST de reforzar territorio:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al reforzar territorio:", resp.StatusCode)
	}

}

func comprobarAcciones(t *testing.T, cookie *http.Cookie) {
	estado := preguntarEstado(t, cookie)

	if len(estado.Acciones) != (logica_juego.NUM_REGIONES + 1) {
		t.Fatal("Se esperaban", logica_juego.NUM_REGIONES, "acciones en el log, y hay", len(estado.Acciones))
	} else {
		t.Log("Contenidos de acciones:", estado.Acciones)
	}

	estado = preguntarEstado(t, cookie)
	if len(estado.Acciones) != 0 {
		t.Fatal("Se esperaban 0 acciones en el log, y hay", len(estado.Acciones))
	} else {
		t.Log("Contenidos de acciones tras leerlas todas:", estado.Acciones)
	}
}

func preguntarEstado(t *testing.T, cookie *http.Cookie) (estado logica_juego.EstadoPartida) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerEstadoPartida", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de preguntar estado:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al preguntar estado:", resp.StatusCode)
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//bodyString := string(body)
	//t.Log("Respuesta al preguntar estado:", bodyString)

	err = json.NewDecoder(resp.Body).Decode(&estado.Acciones)

	return estado
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

func crearUsuario(nombre string, t *testing.T) (cookie *http.Cookie) {
	cookie = nil

	campos := url.Values{
		"nombre":   {nombre},
		"email":    {nombre + "@" + nombre + ".com"},
		"password": {nombre},
	}
	resp, err := http.PostForm("http://localhost:"+os.Getenv(globales.PUERTO_API)+"/registro", campos)
	if err != nil {
		t.Fatal("No se ha podido realizar request POST:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al registrar un usuario:", resp.StatusCode)
	}

	// De ObtenerUsuarioCookie
	for _, c := range resp.Cookies() {
		if c.Name == middleware.NOMBRE_COOKIE_USUARIO { // Es una cookie de usuario
			// Obtener el usuario del valor de la cookie
			nombreCookie := c.Value[:strings.IndexRune(c.Value, middleware.SEPARADOR_VALOR_COOKIE_USUARIO)]
			if nombre != nombreCookie {
				t.Fatal("Obtenido nombre de cookie diferente del esperado:", nombreCookie, "esperaba:", nombre)
			}
			cookie = c
			break
		}
	}

	if cookie == nil {
		t.Fatal("No se ha obtenido una cookie en la respuesta de crear usuario para", nombre)
	}

	return cookie
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
}

func consultarEstadoLobby(cookie *http.Cookie, idPartida int, t *testing.T) (estado vo.EstadoLobby) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerEstadoLobby/"+strconv.Itoa(idPartida), nil)
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

func solicitarAmistad(cookie *http.Cookie, t *testing.T, nombre string) {
	t.Log("Solicitando amistad de userPrincipal a", nombre)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/enviarSolicitudAmistad/"+nombre, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de solicitar amistad:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al solicitar amistad:", resp.StatusCode)
	}
}

func aceptarSolicitudDeAmistad(cookie *http.Cookie, t *testing.T, nombre string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/aceptarSolicitudAmistad/"+nombre, nil) // MAPS :D
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de solicitar amistad:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al aceptar amistad:", resp.StatusCode)
	}
}

func listarAmigos(cookie *http.Cookie, t *testing.T) []string {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/listarAmigos", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de listar amigos:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al listar amigos:", resp.StatusCode)
	} else {
		var amigos vo.ElementoListaNombresUsuario
		err = json.NewDecoder(resp.Body).Decode(&amigos)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al listar amigos:", err)
		}

		t.Log("Respuesta de listarAmigos:", amigos)
		return amigos.Nombres
	}

	return nil
}

func obtenerPerfilUsuario(cookie *http.Cookie, nombre string, t *testing.T) vo.ElementoListaUsuarios {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerPerfil/"+nombre, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de consultar perfil:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al obtener perfil:", resp.StatusCode)
	} else {
		var usuario vo.ElementoListaUsuarios
		err = json.NewDecoder(resp.Body).Decode(&usuario)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al obtener perfil:", err)
		}

		return usuario
	}

	return vo.ElementoListaUsuarios{}
}

func buscarUsuariosSimilares(cookie *http.Cookie, patron string, t *testing.T) []string {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerUsuariosSimilares/"+patron, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de buscar usuarios:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al buscar usuarios:", resp.StatusCode)
	} else {
		var usuarios vo.ElementoListaNombresUsuario
		err = json.NewDecoder(resp.Body).Decode(&usuarios)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al buscar usuarios:", err)
		}
		log.Println("Usuarios recuperados:", usuarios.Nombres)
		return usuarios.Nombres
	}

	return nil
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

func cambiarCartas(t *testing.T, estadoJugador *logica_juego.EstadoJugador, err error, estadoPartida *logica_juego.EstadoPartida, id1, id2, id3, numCanje int) {
	tropasIniciales := estadoJugador.Tropas
	numCartasInicial := len(estadoJugador.Cartas)
	var tropasEsperadas int
	if numCanje < 6 {
		tropasEsperadas = 4 + (numCanje-1)*2
	} else {
		tropasEsperadas = 15 + (numCanje-6)*5
	}
	err = estadoPartida.CambiarCartas("Jugador1", id1, id2, id3)
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
