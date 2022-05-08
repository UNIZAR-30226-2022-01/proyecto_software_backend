package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"testing"
)

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

func TestAbandonoDePartidaEnCurso(t *testing.T) {
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

	unirseAPartida(cookie2, t, 1)
	unirseAPartida(cookie3, t, 1)
	unirseAPartida(cookie4, t, 1)
	unirseAPartida(cookie5, t, 1)
	unirseAPartida(cookie6, t, 1)

	abandonarPartida(cookie, t)
	abandonarPartida(cookie2, t)
	abandonarPartida(cookie3, t)
	abandonarPartida(cookie4, t)
	abandonarPartida(cookie5, t)
	abandonarPartida(cookie6, t)

	_, existe := globales.CachePartidas.ObtenerPartida(1)
	if existe {
		t.Fatal("Se esperaba que la partida dejara de existir en la cache")
	}

	_, err := dao.ObtenerPartida(globales.Db, 1)
	if err == nil {
		t.Fatal("Se esperaba que la partida dejara de existir en la DB")
	}
}
