package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/servidor"
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

// Prueba de unión y abandono de partidas sin errores
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

// Prueba de la funcionalidad de inicio de una partida desde 0
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

	if len(partidaCache.Estado.Acciones) != (int(logica_juego.NUM_REGIONES) + 2) { // 42 regiones y acción de cambio de turno y cambio de fase
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

// Prueba de la funcionalidad de la fase de refuerzo de una partida desde 0
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

	if len(partidaCache.Estado.Acciones) != (int(logica_juego.NUM_REGIONES) + 2) { // 42 regiones y acción de cambio de turno, cambio de fase
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

	saltarTurnos(t, partidaCache, "usuario1")
	// Intentamos cambiar de fase con más de cuatro cartas, se espera error
	partidaCache = comprobarPartidaEnCurso(t, 1)

	numTropasRegionPrevio := partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(numRegion)].NumTropas
	numTropasPrevio := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas

	// Intentamos cambiar de fase con tropas no colocadas, se espera error
	t.Log("Intentamos cambiar de fase con tropas no colocadas, se espera error")
	err := saltarFase(cookie, t)
	partidaCache = comprobarPartidaEnCurso(t, 1)
	t.Log(partidaCache.Estado.EstadosJugadores["usuario1"].Tropas, partidaCache.Estado.Fase)
	if err == nil {
		t.Fatal("Se esperaba error al intentar saltar la fase")
	}
	t.Log("OK: No se ha podido saltar la fase, error:", err)

	// Reforzar con todas las tropas disponibles
	reforzarTerritorio(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas)

	t.Log("Intentamos cambiar de fase con 6 cartas, se espera error")
	partidaCache, err = cambioDeFaseConDemasiadasCartas(t, partidaCache, err, cookie)

	numTropasPost := partidaCache.Estado.EstadosJugadores["usuario1"].Tropas
	numTropasRegionPost := partidaCache.Estado.EstadoMapa[logica_juego.NumRegion(numRegion)].NumTropas

	if numTropasPrevio <= numTropasPost || numTropasPost != 0 {
		t.Fatal("Números de tropas incorrecto al agotarlas en una región. Prev:" + strconv.Itoa(numTropasPrevio) + "Post:" + strconv.Itoa(numTropasPost))
	} else if numTropasRegionPrevio >= numTropasRegionPost || numTropasRegionPost != (numTropasRegionPrevio+numTropasPrevio) {
		t.Fatal("Números de tropas en región incorrecto al agotarlas en una región. Prev:" + strconv.Itoa(numTropasRegionPrevio) + "Post:" + strconv.Itoa(numTropasRegionPost))
	} else {
		t.Log("Tropas asignadas correctamente. Tropas post:", numTropasPost)
	}

	// saltarTurnos(t, partidaCache, "usuario1")
	// Forzar fallo por no tener tropas
	reforzarTerritorioConFallo(t, cookie, numRegion, partidaCache.Estado.EstadosJugadores["usuario1"].Tropas+1)

	// Cambio de fase correcto
	t.Log("Intentamos cambiar de fase con todas las tropas colocadas")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error al saltar de fase:", err)
	}
	t.Log("Se ha saltado la fase correctamente")

	t.Log("Intentamos finalizar el ataque con territorios vacíos, se espera error")
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.Egypt].Ocupante = ""
	globales.CachePartidas.AlmacenarPartida(partidaCache)
	err = saltarFase(cookie, t)
	if err == nil {
		t.Fatal("Se esperaba error al saltar el ataque con territorios vacíos")
	}
	t.Log("OK no se ha podido saltar el ataque con territorios vacíos, error:", err)

	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadoMapa[logica_juego.Egypt].Ocupante = "Jugador1"
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	// Intentamos cambiar de fase con más de cuatro cartas, se espera error
	t.Log("Intentamos cambiar de ataque a fortificación con más de 4 cartas, se espera error")
	partidaCache, err = cambioDeFaseConDemasiadasCartas(t, partidaCache, err, cookie)

	// Cambio de fase correcto
	t.Log("Intentamos cambiar de fase, de ataque a fortificación")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error al saltar de fase:", err)
	}
	t.Log("Se ha saltado la fase correctamente")

	// Cambio de fase correcto
	t.Log("Intentamos cambiar de fase, de fortificación a refuerzo, cediendo el turno")
	err = saltarFase(cookie, t)
	if err != nil {
		t.Fatal("Error al saltar de fase:", err)
	}
	t.Log("Se ha saltado la fase correctamente")

	//saltarTurnos(t, partidaCache, "usuario2")
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

// Prueba de imprimir acciones en JSON
func TestImpresiónEnJSON(t *testing.T) {
	acciones := make([]interface{}, 9)

	acciones[0] = logica_juego.NewAccionRecibirRegion(1, 4, 8, "usuario1")
	acciones[1] = logica_juego.NewAccionInicioTurno("usuario1", 2, 12, 1)
	acciones[2] = logica_juego.NewAccionCambioFase(2, "usuario1")
	acciones[3] = logica_juego.NewAccionCambioCartas(1, 2, true, 2, false)
	acciones[4] = logica_juego.NewAccionReforzar("usuario1", 1, 20)
	acciones[5] = logica_juego.NewAccionAtaque(2, 3, 15, 5, 3, "usuario1", "usuario2")
	acciones[6] = logica_juego.NewAccionOcupar(2, 3, 10, 5, "usuario1", "usuario2")
	acciones[7] = logica_juego.NewAccionFortificar(7, 9, 10, 8, "usuario1")
	acciones[8] = logica_juego.NewAccionObtenerCarta(logica_juego.Carta{Tipo: logica_juego.Infanteria, Region: logica_juego.Egypt}, "usuario1")

	for _, a := range acciones {
		serializarAJSONEImprimir(t, a)
	}
}

// Prueba de consistencia entre la caché de partidas y la base de datos
func TestConsistencia(t *testing.T) {
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
	comprobarConsistenciaEnCurso(t, partidaCache)

	t.Log("Uniéndose a partida...")
	unirseAPartida(cookie2, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)
	comprobarConsistenciaEnCurso(t, partidaCache)

	unirseAPartida(cookie3, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)
	comprobarConsistenciaEnCurso(t, partidaCache)

	unirseAPartida(cookie4, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)
	comprobarConsistenciaEnCurso(t, partidaCache)

	unirseAPartida(cookie5, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 1 {
		t.Fatal("No hay partidas, aunque aún no ha empezado.")
	}
	partidaCache = comprobarPartidaNoEnCurso(t, 1)
	comprobarConsistenciaEnCurso(t, partidaCache)

	unirseAPartida(cookie6, t, 1)
	partidas = obtenerPartidas(cookie, t)
	if len(partidas) != 0 {
		t.Fatal("Hay partidas, aunque ya ha tenido que empezar:", partidas)
	}
	partidaCache = comprobarPartidaEnCurso(t, 1)
	comprobarConsistenciaEnCurso(t, partidaCache)

	comprobarConsistenciaAcciones(t, partidaCache)
}

func TestBaraja(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	estadoPartida := logica_juego.CrearEstadoPartida([]string{"Jugador1", "Jugador2", "Jugador3"})
	estadoPartida.TurnoJugador = 0
	estadoPartida.Fase = logica_juego.Fortificar
	robarBarajaCompleta(&estadoPartida, t)
	if len(estadoPartida.Cartas) > 0 {
		t.Fatal("La baraja debería estar vacía, pero tiene", len(estadoPartida.Cartas), "cartas")
	}
	estadoJugador := estadoPartida.EstadosJugadores["Jugador1"]
	t.Log("Baraja completa:", estadoJugador.Cartas)
	var err error
	var carta logica_juego.Carta
	var baraja []logica_juego.Carta

	// Prueba función retirar carta por ID
	// Retiramos todas las cartas de la mano del jugador
	for i := 0; i < 44; i++ {
		carta, estadoJugador.Cartas, err = logica_juego.RetirarCartaPorID(i, estadoJugador.Cartas)
		if err != nil {
			t.Fatal("Error al retirar cartas del jugador:", err)
		}
		if carta.IdCarta != i {
			t.Fatal("La carta recuperada no es correcta, id:", carta.IdCarta, "se esperaba: ", i)
		}
		baraja = append(baraja, carta)
	}
	t.Log("Cartas en la mano del jugador:", estadoJugador.Cartas)
	if len(estadoJugador.Cartas) > 0 {
		t.Fatal("El jugador no debería tener cartas")
	}

	// Compiamos las cartas al montón de descartes
	estadoPartida.Descartes = baraja

	// Tomamos una carta de la baraja
	// Al estar vacía, toma los descartes y los rebaraja
	estadoPartida.HaConquistado = true
	estadoPartida.HaRecibidoCarta = false

	err = estadoPartida.RecibirCarta("Jugador1")
	if err != nil {
		t.Fatal("Error al recibir carta:", err)
	}
	estadoPartida.HaConquistado = true
	estadoPartida.HaRecibidoCarta = false

	// La baraja debería tener 43 cartas
	if len(estadoPartida.Cartas) != 43 {
		t.Fatal("La baraja debería tener 43 cartas, pero tiene:", len(estadoPartida.Cartas))
	}

	// Volvemos a tomar toda la baraja
	robarBarajaCompleta(&estadoPartida, t)
	if len(estadoPartida.Cartas) > 0 {
		t.Fatal("La baraja debería estar vacía, pero tiene:", len(estadoPartida.Cartas), "cartas")
	}

	// Prueba de canjes
	// Cambiamos 3 cartas de infantería
	estadoPartida.Fase = logica_juego.Refuerzo
	t.Log("Cambiando 3 cartas de infanteria")
	cambiarCartas(t, estadoJugador, &estadoPartida, 0, 1, 2, 1)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos 3 cartas de caballeria
	t.Log("Cambiando 3 cartas de caballeria")
	cambiarCartas(t, estadoJugador, &estadoPartida, 18, 19, 20, 2)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos 3 cartas de artilleria
	t.Log("Cambiando 3 cartas de artilleria")
	cambiarCartas(t, estadoJugador, &estadoPartida, 36, 37, 38, 3)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos una de cada
	t.Log("Cambiando 3 cartas, una de cada tipo")
	cambiarCartas(t, estadoJugador, &estadoPartida, 3, 21, 39, 4)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Cambiamos 2 cartas de un tipo + un comodín
	t.Log("Cambiando 3 cartas, una de cada tipo")
	cambiarCartas(t, estadoJugador, &estadoPartida, 4, 5, 42, 5)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	// Prueba de errores en el canje
	// Cambiar cartas que no tenemos
	t.Log("Intentamos cambiar cartas que no tenemos, se espera error")
	err = estadoPartida.CambiarCartas("Jugador1", 0, 1, 2)
	if err == nil {
		t.Fatal("Se esperaba obtener error al cambiar con cartas que el jugador no tiene")
	}

	t.Log("OK: No se ha podido cambiar con cartas que el jugador no tiene, error:", err)

	// Cambiar con cartas de distinto tipo
	t.Log("Intentamos cambiar con cartas de distinto tipo, se espera error")
	err = estadoPartida.CambiarCartas("Jugador1", 6, 7, 23)
	if err == nil {
		t.Fatal("Se esperaba obtener error al cambiar con cartas de distintos tipos")
	}

	t.Log("OK: No se ha podido cambiar con cartas de diferentes tipos, error:", err)

	// Probamos bonificación por territorio
	// TODO probar caso en el que más de una carta recibe bonificación por territorio
	carta, estadoJugador.Cartas, err = logica_juego.RetirarCartaPorID(10, estadoJugador.Cartas)
	if err != nil {
		t.Fatal("Error al tomar la carta de la mano del jugador")
	}
	estadoJugador.Cartas = append(estadoJugador.Cartas, carta)
	region := carta.Region
	estadoRegion := estadoPartida.EstadoMapa[region]
	estadoRegion.Ocupante = "Jugador1"
	tropasIniciales := estadoRegion.NumTropas
	t.Log("Cambiando 3 cartas, una de ellas con bonificación por territorio")
	cambiarCartas(t, estadoJugador, &estadoPartida, 10, 11, 12, 6)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)
	if estadoRegion.NumTropas-tropasIniciales != 2 {
		t.Fatal("No se han recibido tropas adicionales")
	}
	t.Log("Inicialmente, la región tenía", tropasIniciales, "tropas, tras el canje hay un total de ", estadoRegion.NumTropas-tropasIniciales)

	// Comprobar que se añade correctamente la acción de cambio de cartas
	accion := estadoPartida.Acciones[len(estadoPartida.Acciones)-1]
	accionCambio, ok := accion.(logica_juego.AccionCambioCartas)
	if !ok {
		t.Fatal("La última acción no es un cambio de cartas")
	}

	if !accionCambio.ObligadoAHacerCambios {
		t.Fatal("El jugador debería estar obligado a cambiar, al tener más de 4 cartas")
	}

	if !accionCambio.BonificacionObtenida {
		t.Fatal("No se ha obtenido bonificación por territorio al cambiar las cartas")
	}

	t.Log("El último cambio de cartas fue:", accionCambio)

	// Probar cambios de cartas de número > 6
	t.Log("Cambiando 3 cartas, 7º cambio")
	cambiarCartas(t, estadoJugador, &estadoPartida, 30, 31, 32, 7)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)

	t.Log("Cambiando 3 cartas, 8º cambio")
	cambiarCartas(t, estadoJugador, &estadoPartida, 33, 34, 35, 8)
	invarianteNumeroDeCartas(estadoPartida, *estadoJugador, t)
}

func TestAsignacionTropas(t *testing.T) {
	// TODO modificar prueba para usar cambio de turno y evitar que AsignarTropasRefuerzo sea pública
	// Casos a probar
	// Territorios < 12
	// Territorios >= 12
	// Ocupa algún continente

	t.Log("Purgando DB...")
	purgarDB()

	estadoPartida := logica_juego.CrearEstadoPartida([]string{"Jugador1", "Jugador2", "Jugador3", "Jugador4", "Jugador5", "Jugador6"})
	estadoPartida.TurnoJugador = 5
	estadoPartida.Fase = logica_juego.Fortificar
	estadoJugador := estadoPartida.EstadosJugadores["Jugador1"]

	// Desocupamos el mapa
	for region, _ := range estadoPartida.EstadoMapa {
		estadoPartida.EstadoMapa[region].Ocupante = ""
	}

	estadoJugador.Tropas = 0

	// Comprobamos que se asignan 3 ejércitos en caso de tener menos de 12 territorios
	t.Log("El jugador no ocupa ningún territorio")
	estadoPartida.AsignarTropasRefuerzo("Jugador1")
	if estadoJugador.Tropas != 3 {
		t.Fatal("El jugador debería tener 3 tropas pero tiene", estadoJugador.Tropas)
	}
	t.Log("El jugador ha recibido", estadoJugador.Tropas, "tropas al principio del turno")

	// El jugador ocupa Asia
	// Deberá recibir 11 ejércitos
	t.Log("El jugador ocupa Asia (7), con 12 territorios (4)")
	for _, region := range logica_juego.Continentes["Asia"].Regiones {
		estadoPartida.EstadoMapa[region].Ocupante = "Jugador1"
	}

	estadoJugador.Tropas = 0
	estadoPartida.AsignarTropasRefuerzo("Jugador1")
	if estadoJugador.Tropas != 11 {
		t.Fatal("El jugador debería tener 11 tropas pero tiene", estadoJugador.Tropas)
	}
	t.Log("El jugador ha recibido", estadoJugador.Tropas, "tropas al principio del turno")

	// El jugador ocupa Asia y Europa
	// Un total de 19 regiones
	// Deberá recibir 5 ejércitos por Europa, 7 por Asia y 6 por territorios ocupados
	// Total ejércitos = 18
	t.Log("El jugador ocupa Asia (7) y Europa (5), con 19 territorios (6)")
	for _, region := range logica_juego.Continentes["Europa"].Regiones {
		estadoPartida.EstadoMapa[region].Ocupante = "Jugador1"
	}

	estadoJugador.Tropas = 0
	estadoPartida.AsignarTropasRefuerzo("Jugador1")
	if estadoJugador.Tropas != 18 {
		t.Fatal("El jugador debería tener 18 tropas pero tiene", estadoJugador.Tropas)
	}
	t.Log("El jugador ha recibido", estadoJugador.Tropas, "tropas al principio del turno")
}
