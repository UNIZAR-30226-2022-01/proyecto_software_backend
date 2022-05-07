package integracion

import (
	"log"
	"net/http"
	"testing"
)

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

	// Comprobamos que no se pueden enviar solicitudes de amistad a un usuario, en caso de que él nos haya enviado una
	t.Log("Enviamos una solicitud de amistad a un usuario que ya nos la había solicitado, se espera error")
	solicitarAmistadConError(cookie, t, "Amigo1")

	solicitudesPendientes := consultarSolicitudesPendientes(cookie, t)
	if len(solicitudesPendientes) != len(amigos) {
		t.Fatal("No se han recuperado todas las solicitudes pendientes")
	}

	for i := range amigos {
		if amigos[i] != solicitudesPendientes[i] {
			t.Fatal("No se han recuperado todas las solicitudes pendientes")
		}
	}

	// Rechazamos todas las solicitudes menos la última
	for _, a := range amigos[:len(amigos)-1] {
		rechazarSolicitudDeAmistad(cookie, t, a)
	}

	solicitudesPendientes = consultarSolicitudesPendientes(cookie, t)
	if len(solicitudesPendientes) != 1 {
		t.Fatal("Hay más de una solicitud pendiente")
	}

	// Solicita amistad al resto de usuarios
	for _, a := range amigos[:len(amigos)-1] {
		solicitarAmistad(cookie, t, a)
	}

	// Comprobamos que las amistades pendientes no aparecen como amigos
	listaAmigos := listarAmigos(cookie, t)
	t.Log("Lista de amigos, debería estar vacía:", listaAmigos)
	if len(listaAmigos) > 0 {
		t.Fatal("Las solicitudes de amistad pendientes no deberían aparecer en la lista de amigos")
	}

	// Cada uno acepta la solicitud
	for _, c := range cookiesAmigos[:len(cookiesAmigos)-1] {
		aceptarSolicitudDeAmistad(c, t, "usuario")
	}

	aceptarSolicitudDeAmistad(cookie, t, "Amigo5")

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

	// Intentamos modificar la biografía del usuario
	biografia := "Mejor jugador del Risk 2021"
	t.Log("Intentamos modificar la biografía del usuario")
	modificarBiografia(cookie, biografia, t)
	usuario := obtenerPerfilUsuario(cookie, "usuario", t)
	if usuario.Biografia != biografia {
		t.Fatal("No se ha cambiado la biografía correctamente")
	}
	t.Log("La nueva biografía es:", usuario.Biografia)

	// Buscamos usuarios cuyo nombre empiece por "Amigo"
	resultadoBusqueda := buscarUsuariosSimilares(cookie, "Amigo", t)
	log.Println("amigos:", resultadoBusqueda)
	if len(amigos) != len(resultadoBusqueda) {
		t.Fatal("No se han recuperado todos los usuarios con nombre empezado por Amigo")
	}

	// Probamos que se recuperen los mismos amigos aunque no coincidan mayúsculas y minúsculas
	resultadoBusqueda = buscarUsuariosSimilares(cookie, "aMIgo", t)
	log.Println("amigos:", resultadoBusqueda)
	if len(amigos) != len(resultadoBusqueda) {
		t.Fatal("No se han recuperado todos los usuarios con nombre empezado por Amigo al variar mayúsculas y minúsuclas")
	}

	for i := range amigos {
		if amigos[i] != resultadoBusqueda[i].Nombre {
			t.Fatal("No se han recuperado todos los usuarios con nombre empezado por Amigo")
		} else if !resultadoBusqueda[i].EsAmigo {
			t.Fatal("Uno de los amigos se ha devuelto como no amigo:", resultadoBusqueda[i])
		}
	}

	// Crea varios usuarios no amigos y se comprueba que no aparecen como amigos
	noAmigos := []string{"NoAmigo1", "NoAmigo2", "NoAmigo3"}
	cookiesNoAmigos := make([]*http.Cookie, 5)
	for i, a := range noAmigos {
		cookiesNoAmigos[i] = crearUsuario(a, t)
	}

	resultadoBusquedaNoAmigos := buscarUsuariosSimilares(cookie, "NoAmi", t)
	log.Println("noAmigos:", resultadoBusquedaNoAmigos)
	if len(noAmigos) != len(resultadoBusquedaNoAmigos) {
		t.Fatal("No se han recuperado todos los usuarios con nombre empezado por NoAmi")
	}

	for i := range noAmigos {
		if noAmigos[i] != resultadoBusquedaNoAmigos[i].Nombre {
			t.Fatal("No se han recuperado todos los usuarios con nombre empezado por NoAmi")
		} else if resultadoBusquedaNoAmigos[i].EsAmigo {
			t.Fatal("Uno de los (no) amigos se ha devuelto como amigo:", resultadoBusquedaNoAmigos[i])
		}
	}

	// Prueba de bandera de amigo al consultar usuarios individuales
	amigo := obtenerPerfilUsuario(cookie, "Amigo1", t)
	if !amigo.EsAmigo {
		t.Fatal("Amigo1 se ha devuelto como no amigo al consultarlo:", amigo)
	}

	noAmigo := obtenerPerfilUsuario(cookie, "NoAmigo1", t)
	if noAmigo.EsAmigo {
		t.Fatal("NoAmigo1 se ha devuelto como amigo al consultarlo:", noAmigo)
	}

	// Borrar amistades
	t.Log("Intentamos borrar un amigo")
	eliminarAmigo(cookie, "Amigo5", t)

	listaAmigos = listarAmigos(cookie, t)
	if len(listaAmigos) != 4 {
		t.Fatal("No se ha borrado el amigo correctamente")
	}
	t.Log("OK, se ha borrado el amigo")

	t.Log("Intentamos borrar a otro amigo")
	eliminarAmigo(cookie, "Amigo1", t)

	listaAmigos = listarAmigos(cookie, t)
	if len(listaAmigos) != 3 {
		t.Fatal("No se ha borrado el amigo correctamente")
	}
	t.Log("OK, se ha borrado el amigo")
}
