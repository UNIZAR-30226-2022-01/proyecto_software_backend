package integracion_test

import (
	"backend/globales"
	"backend/middleware"
	"backend/servidor"
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
	"time"
)

// Función que se ejecuta automáticamente antes de los test
func init() {
	// Inyecta las variables de entorno
	os.Setenv("DIRECCION_DB", "postgres")
	os.Setenv("DIRECCION_DB_TESTS", "localhost")
	os.Setenv("PUERTO_API", "8090")
	os.Setenv("PUERTO_WEB", "8080")
	os.Setenv("USUARIO_DB", "postgres")
	os.Setenv("PASSWORD_DB", "postgres")

	go servidor.IniciarServidor(true)
	time.Sleep(5 * time.Second)
}

// Prueba, del lado del cliente, de:
//		Crear una serie de usuarios
//		Crear una serie de partidas
//		Obtener y comprobar ordenación de partidas
//
// Asume una BD limpia
func TestCreacionYObtencionPartidas(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuarios...")
	cookiesCreadores := make([]*http.Cookie, 6)
	cookiesCreadores[0] = crearUsuario("creadorP1", t)
	cookiesCreadores[1] = crearUsuario("creadorP2", t)
	cookiesCreadores[2] = crearUsuario("creadorP3", t)
	cookiesCreadores[3] = crearUsuario("creadorP4", t)
	cookiesCreadores[4] = crearUsuario("creadorP5", t)
	cookiesCreadores[5] = crearUsuario("creadorP6", t)

	cookieUsuarioPrincipal := crearUsuario("userPrincipal", t)
	t.Log("Cookie a usar:", cookieUsuarioPrincipal)

	cookiesAmigos := make([]*http.Cookie, 5)
	cookiesAmigos[0] = crearUsuario("amigo1", t)
	cookiesAmigos[1] = crearUsuario("amigo2", t)
	cookiesAmigos[2] = crearUsuario("amigo3", t)
	cookiesAmigos[3] = crearUsuario("amigo4", t)
	cookiesAmigos[4] = crearUsuario("amigo5", t)

	cookiesNoAmigos := make([]*http.Cookie, 8)
	cookiesNoAmigos[0] = crearUsuario("NoAmigo1", t)
	cookiesNoAmigos[1] = crearUsuario("NoAmigo2", t)
	cookiesNoAmigos[2] = crearUsuario("NoAmigo3", t)
	cookiesNoAmigos[3] = crearUsuario("NoAmigo4", t)
	cookiesNoAmigos[4] = crearUsuario("NoAmigo5", t)
	cookiesNoAmigos[5] = crearUsuario("NoAmigo6", t)
	cookiesNoAmigos[6] = crearUsuario("NoAmigo7", t)
	cookiesNoAmigos[7] = crearUsuario("NoAmigo8", t)

	// 3 privadas
	crearPartida(cookiesCreadores[0], t, false)
	crearPartida(cookiesCreadores[1], t, false)
	crearPartida(cookiesCreadores[2], t, false)

	// 3 públicas
	crearPartida(cookiesCreadores[3], t, true)
	crearPartida(cookiesCreadores[4], t, true)
	crearPartida(cookiesCreadores[5], t, true)

	for i, _ := range cookiesAmigos {
		solicitarAmistad(cookieUsuarioPrincipal, t, "amigo"+strconv.Itoa(i+1))
	}

	for _, c := range cookiesAmigos {
		aceptarSolicitudDeAmistad(c, t, "userPrincipal")
	}

	// P1 privada con 2 amigos, 1 no
	unirseAPartida(cookiesAmigos[0], t, 1)
	unirseAPartida(cookiesAmigos[1], t, 1)
	unirseAPartida(cookiesNoAmigos[0], t, 1)

	// P2 privada con 1 amigo, 2 no
	unirseAPartida(cookiesAmigos[2], t, 2)
	unirseAPartida(cookiesNoAmigos[1], t, 2)
	unirseAPartida(cookiesNoAmigos[2], t, 2)

	// P3 privada con 0 amigos, 1 no
	unirseAPartida(cookiesNoAmigos[3], t, 3)

	// P4 pública con 2 amigos, 1 no
	unirseAPartida(cookiesAmigos[3], t, 4)
	unirseAPartida(cookiesAmigos[4], t, 4)
	unirseAPartida(cookiesNoAmigos[4], t, 4)

	// P5 pública con 0 amigos, 2 no
	unirseAPartida(cookiesNoAmigos[5], t, 5)
	unirseAPartida(cookiesNoAmigos[6], t, 5)

	// P6 pública con 0 amigos, 1 no
	unirseAPartida(cookiesNoAmigos[7], t, 6)

	// Orden: P1, P2, P3, P4, P5, P6

	partidas := obtenerPartidas(cookieUsuarioPrincipal, t)

	if partidas[0].IdPartida != 1 {
		t.Fatal("Se esperaba obtener partida de ID", 1, "en posición", 0)
	} else if partidas[1].IdPartida != 2 {
		t.Fatal("Se esperaba obtener partida de ID", 2, "en posición", 1)
	} else if partidas[2].IdPartida != 4 {
		t.Fatal("Se esperaba obtener partida de ID", 4, "en posición", 2)
	} else if partidas[3].IdPartida != 5 {
		t.Fatal("Se esperaba obtener partida de ID", 5, "en posición", 3)
	} else if partidas[4].IdPartida != 6 {
		t.Fatal("Se esperaba obtener partida de ID", 1, "en posición", 4)
	} else if partidas[5].IdPartida != 3 {
		t.Fatal("Se esperaba obtener partida de ID", 1, "en posición", 5)
	} else {
		t.Log("Partidas ordenadas correctamente!")
	}
}

func TestUnionYAbandonoDePartidas(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuario...")
	cookie := crearUsuario("usuario1", t)
	cookie2 := crearUsuario("usuario2", t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, false)

	partidas := obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas.")
	}

	unirseAPartida(cookie2, t, 1)
	abandonarLobby(cookie, t)
	partidas = obtenerPartidas(cookie, t)
	// Aunque se ha ido el creador, el lobby debería seguir existiendo al haber otro usuario
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque hay un usuario aún.")
	}

	abandonarLobby(cookie2, t)
	// Ahora se de debería haber borrado
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 0 {
		t.Fatal("Sigue habiendo partidas tras quedar con 0 usuarios.")
	}
}

func TestInicioPartida(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuarios...")
	cookie := crearUsuario("usuario1", t)
	cookie2 := crearUsuario("usuario2", t)
	cookie3 := crearUsuario("usuario3", t)
	cookie4 := crearUsuario("usuario4", t)
	cookie5 := crearUsuario("usuario5", t)
	cookie6 := crearUsuario("usuario6", t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, true)

	partidas := obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas.")
	}

	partidaCache := comprobarPartidaNoEnCurso(t, 1)

	t.Log("Uniéndose a partida...")
	unirseAPartida(cookie2, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie3, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie4, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie5, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie6, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 0 {
		t.Fatal("Hay partidas, aunque ya ha tenido que empezar:", partidas)
	}
	partidaCache = comprobarPartidaEnCurso(t, 1)

	if len(partidaCache.Estado.Acciones) != (int(vo.NUM_REGIONES) + 1) { // 42 regiones y acción de cambio de turno
		t.Fatal("No se han asignado todas las regiones. Regiones asignadas:", len(partidaCache.Estado.Acciones))
	}

	for _, jugador := range partidaCache.Estado.Jugadores {
		t.Log("Estado de ", jugador, ":", partidaCache.Estado.EstadosJugadores[jugador])
	}

	comprobarAcciones(t, cookie)
	comprobarAcciones(t, cookie2)
	comprobarAcciones(t, cookie3)
	comprobarAcciones(t, cookie4)
	comprobarAcciones(t, cookie5)
	comprobarAcciones(t, cookie6)
}

func TestFaseRefuerzoInicial(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuarios...")
	cookie := crearUsuario("usuario1", t)
	cookie2 := crearUsuario("usuario2", t)
	cookie3 := crearUsuario("usuario3", t)
	cookie4 := crearUsuario("usuario4", t)
	cookie5 := crearUsuario("usuario5", t)
	cookie6 := crearUsuario("usuario6", t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, true)

	partidas := obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas.")
	}

	partidaCache := comprobarPartidaNoEnCurso(t, 1)

	t.Log("Uniéndose a partida...")
	unirseAPartida(cookie2, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie3, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie4, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie5, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)

	unirseAPartida(cookie6, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 0 {
		t.Fatal("Hay partidas, aunque ya ha tenido que empezar:", partidas)
	}
	partidaCache = comprobarPartidaEnCurso(t, 1)

	if len(partidaCache.Estado.Acciones) != (int(vo.NUM_REGIONES) + 1) { // 42 regiones y acción de cambio de turno
		t.Fatal("No se han asignado todas las regiones. Regiones asignadas:", len(partidaCache.Estado.Acciones))
	}

	numRegion := 0
	// Buscar región ocupada por usuario1
	for i := vo.Eastern_australia; i <= vo.Alberta; i++ {
		if partidaCache.Estado.EstadoMapa[i].Ocupante == "usuario1" {
			numRegion = int(i)
			break
		}
	}

	numTropasRegionPrevio := partidaCache.Estado.EstadoMapa[vo.NumRegion(numRegion)].NumTropas
	numTropasPrevio := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas

	saltarTurnos(t, partidaCache, "usuario1")

	// Reforzar con todas las tropas disponibles
	reforzarTerritorio(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas)

	partidaCache = comprobarPartidaEnCurso(t, 1)

	numTropasPost := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas
	numTropasRegionPost := partidaCache.Estado.EstadoMapa[vo.NumRegion(numRegion)].NumTropas

	if numTropasPrevio <= numTropasPost || numTropasPost != 0 {
		t.Fatal("Números de tropas incorrecto al agotarlas en una región. Prev:" + strconv.Itoa(numTropasPrevio) + "Post:" + strconv.Itoa(numTropasPost))
	} else if numTropasRegionPrevio >= numTropasRegionPost || numTropasRegionPost != (numTropasRegionPrevio+numTropasPrevio) {
		t.Fatal("Números de tropas en región incorrecto al agotarlas en una región. Prev:" + strconv.Itoa(numTropasRegionPrevio) + "Post:" + strconv.Itoa(numTropasRegionPost))
	} else {
		t.Log("Tropas asignadas correctamente. Tropas post:", numTropasPost)
	}

	saltarTurnos(t, partidaCache, "usuario1")
	// Forzar fallo por no tener tropas
	reforzarTerritorioConFallo(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas+1)
	saltarTurnos(t, partidaCache, "usuario2")
	// Forzar fallo por estar fuera de turno
	reforzarTerritorioConFallo(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	t.Log("Acciones al final:", partidaCache.Estado.Acciones)
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

	if len(estado.Acciones) != (vo.NUM_REGIONES + 1) {
		t.Fatal("Se esperaban", vo.NUM_REGIONES, "acciones en el log, y hay", len(estado.Acciones))
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

func preguntarEstado(t *testing.T, cookie *http.Cookie) (estado vo.EstadoPartida) {
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
