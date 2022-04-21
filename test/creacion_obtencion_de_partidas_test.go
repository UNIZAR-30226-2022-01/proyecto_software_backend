package integracion

import (
	"net/http"
	"strconv"
	"testing"
)

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
	estadoLobby := consultarEstadoLobby(cookiesAmigos[0], t)
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
	estadoLobby = consultarEstadoLobby(cookiesAmigos[3], t)
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
