package integracion

import (
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
