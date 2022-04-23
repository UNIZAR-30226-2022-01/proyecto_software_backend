package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"testing"
	"time"
)

func TestEliminacionPartidas(t *testing.T) {
	t.Skip()
	// Ejecutar test manualmente modificando el temporizador de limpiarCache en globales.go

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

	time.Sleep(3 * 6 * time.Second) // Asumiendo tiempos de expulsión y pasadas de limpieza de 2 segundos
	_, err := dao.ObtenerPartida(globales.Db, 1)

	if err == nil {
		t.Log("La partida no debería aparecer en la base de datos por tener todos los jugadores expulsados")
	}

	_, existe := globales.CachePartidas.ObtenerPartida(1)

	if existe {
		t.Log("La partida no debería aparecer en la cache por tener todos los jugadores expulsados")
	}
}
