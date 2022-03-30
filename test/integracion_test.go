package integracion

import (
	"backend/logica_juego"
	"backend/servidor"
	"net/http"
	"os"
	"strconv"
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

	// Consultamos el estado del lobby
	estadoLobby := consultarEstadoLobby(cookiesAmigos[0], 1, t)
	if estadoLobby.EsPublico {
		t.Fatal("La partida debería ser privada")
	}
	if estadoLobby.Jugadores != 4 {
		t.Fatal("Debería haber 4 jugadores en el lobby")
	}
	if estadoLobby.MaxJugadores != 6 {
		t.Fatal("El máximo de jugadores debería ser 6")
	}
	if estadoLobby.EnCurso {
		t.Fatal("La partida no debería estar en curso")
	}

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

	// Consultamos el estado del lobby
	estadoLobby = consultarEstadoLobby(cookiesAmigos[0], 4, t)
	if !estadoLobby.EsPublico {
		t.Fatal("La partida debería ser pública")
	}
	if estadoLobby.Jugadores != 4 {
		t.Fatal("Debería haber 4 jugadores en el lobby")
	}
	if estadoLobby.MaxJugadores != 6 {
		t.Fatal("El máximo de jugadores debería ser 6")
	}
	if estadoLobby.EnCurso {
		t.Fatal("La partida no debería estar en curso")
	}

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

	if len(partidaCache.Estado.Acciones) != (int(logica_juego.NUM_REGIONES) + 1) { // 42 regiones y acción de cambio de turno
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

	if len(partidaCache.Estado.Acciones) != (int(logica_juego.NUM_REGIONES) + 1) { // 42 regiones y acción de cambio de turno
		t.Fatal("No se han asignado todas las regiones. Regiones asignadas:", len(partidaCache.Estado.Acciones))
	}

	numRegion := 0
	// Buscar región ocupada por usuario1
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		if partidaCache.Estado.EstadoMapa[i].Ocupante == "usuario1" {
			numRegion = int(i)
			break
		}
	}

	numTropasRegionPrevio := partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(numRegion)].NumTropas
	numTropasPrevio := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas

	saltarTurnos(t, partidaCache, "usuario1")

	// Reforzar con todas las tropas disponibles
	reforzarTerritorio(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas)

	partidaCache = comprobarPartidaEnCurso(t, 1)

	numTropasPost := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas
	numTropasRegionPost := partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(numRegion)].NumTropas

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

	// Prueba para la consulta de amigos pendientes
	// El resto de usuarios solicitan amistad al primer usuario
	for _, c := range cookiesAmigos {
		solicitarAmistad(c, t, "usuario")
	}

	solicitudesPendientes := consultarSolicitudesPendientes(cookie, t)
	if len(solicitudesPendientes) != len(amigos) {
		t.Fatal("No se han recuperado todas las solicitudes pendientes")
	}

	for i := range amigos {
		if amigos[i] != solicitudesPendientes[i] {
			t.Fatal("No se han recuperado todas las solicitudes pendientes")
		}
	}

	// Rechazamos todas las solicitudes
	for _, a := range amigos {
		rechazarSolicitudDeAmistad(cookie, t, a)
	}

	solicitudesPendientes = consultarSolicitudesPendientes(cookie, t)
	if len(solicitudesPendientes) != 0 {
		t.Fatal("Se han recuperado solicitudes pendientes cuando no debería haberlas")
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

	// Recuperamos la información de perfil del primer usuario
	usuarioRecuperado := obtenerPerfilUsuario(cookie, "usuario", t)
	if usuarioRecuperado.NombreUsuario != "usuario" {
		t.Fatal("No se ha obtenido correctamente el perfil del usuario")
	}

	if usuarioRecuperado.Email != "usuario@usuario.com" {
		t.Fatal("No se ha obtenido correctamente el perfil del usuario")
	}

	// TODO -> probar si recuperamos la biografia y otros campos correctamente una vez se puedan modificar

	// Buscamos usuarios cuyo nombre empiece por "Amigo"
	resultadoBusqueda := buscarUsuariosSimilares(cookie, "Amigo", t)
	if len(amigos) != len(resultadoBusqueda) {
		t.Fatal("No se han recuperado todos los usuarios con nombre empezado por Amigo")
	}

	for i := range amigos {
		if amigos[i] != resultadoBusqueda[i] {
			t.Fatal("No se han recuperado todos los usuarios con nombre empezado por Amigo")
		}
	}

}
