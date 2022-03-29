package integracion_test

import (
	"backend/globales"
	"backend/middleware"
	"backend/servidor"
	"backend/vo"
	"encoding/json"
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

	for i := range cookiesAmigos {
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

// Prueba las llamadas a la API de listar amigos, obtener información de perfil y buscar usuarios que coincidan con
// un nombre
func TestFuncionesSociales(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	cookie := crearUsuario("usuario", t)
	amigos := []string{"Amigo1", "Amigo2", "Amigo3", "Amigo4", "Amigo5"}
	cookiesAmigos := make([]*http.Cookie, 5)
	for i, a := range amigos {
		cookiesAmigos[i] = crearUsuario(a, t)
	}

	// Solicita amistad al resto de usuarios
	for _, a := range amigos {
		solicitarAmistad(cookie, t, a)
	}

	// Cada uno acepta la solicitud
	for _, c := range cookiesAmigos {
		aceptarSolicitudDeAmistad(c, t, "usuario")
	}

	amigosRegistrados := listarAmigos(cookie, t)
	if len(amigos) != len(amigosRegistrados) {
		t.Fatal("No se han recuperado todos los amigos")
	}

	for i := range amigos {
		if amigos[i] != amigosRegistrados[i] {
			t.Fatal("No se han recuperado todos los amigos")
		}
	}
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
		var amigos vo.ElementoListaAmigos
		err = json.NewDecoder(resp.Body).Decode(&amigos)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al listar amigos:", err)
		}

		t.Log("Respuesta de listarAmigos:", amigos)
		return amigos.Nombres
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
