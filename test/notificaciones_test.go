package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"testing"
)

func TestNotificaciones(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	// Creación e inicio de la partida

	t.Log("Creando usuarios...")
	cookie := crearUsuario("usuario1", t)
	cookie2 := crearUsuario("usuario2", t)
	cookie3 := crearUsuario("usuario3", t)
	cookie4 := crearUsuario("usuario4", t)
	cookie5 := crearUsuario("usuario5", t)
	cookie6 := crearUsuario("usuario6", t)

	t.Log("Creando partida...")
	crearPartida(cookie, t, true)

	t.Log("Uniéndose a partida...")
	unirseAPartida(cookie2, t, 1)
	unirseAPartida(cookie3, t, 1)
	unirseAPartida(cookie4, t, 1)
	unirseAPartida(cookie5, t, 1)
	unirseAPartida(cookie6, t, 1)

	partidaCache := comprobarPartidaEnCurso(t, 1)

	saltarTurnos(t, partidaCache, "usuario1")
	notificaciones := obtenerNotificaciones(t, cookie)

	if len(notificaciones) != 1 {
		t.Fatal("Se esperaba una notificación de turno, obtenido:", notificaciones)
	} else {
		t.Log("Notificaciones tras tener turno pendiente:", notificaciones)
	}

	solicitarAmistad(cookie2, t, "usuario1")
	notificaciones = obtenerNotificaciones(t, cookie)

	if len(notificaciones) != 2 {
		t.Fatal("Se esperaba una notificación de turno y otra de amistad, obtenido:", notificaciones)
	} else {
		t.Log("Notificaciones tras tener turno pendiente y amistad pendiente:", notificaciones)
	}

	// Probar notificaciones de puntos
	dao.OtorgarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario1"}, 100, true)
	notificaciones = obtenerNotificaciones(t, cookie)
	if len(notificaciones) != 3 {
		t.Fatal("Se esperaba una notificación de turno, otra de amistad y otra de puntos, obtenido:", notificaciones)
	} else {
		t.Log("Notificaciones tras tener turno pendiente, amistad pendiente y otra de puntos:", notificaciones)
	}

	notificaciones = obtenerNotificaciones(t, cookie)
	if len(notificaciones) != 2 {
		t.Fatal("Se esperaba que se borrara la notificación de puntos, obtenido:", notificaciones)
	} else {
		t.Log("Notificaciones tras tener turno pendiente y amistad pendiente y habiendo obtenido ya una de puntos:", notificaciones)
	}
}
