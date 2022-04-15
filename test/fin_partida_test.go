package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"net/http"
	"testing"
)

func TestFinPartida(t *testing.T) {
	t.Log("Purgando DB...")
	purgarDB()

	t.Log("Creando usuarios...")
	cookies := make([]*http.Cookie, 3)
	cookies[0] = crearUsuario("usuario1", t)
	cookies[1] = crearUsuario("usuario2", t)
	cookies[2] = crearUsuario("usuario3", t)

	t.Log("Creando partida...")
	crearPartidaReducida(cookies[0], t, true)
	unirseAPartida(cookies[1], t, 1)
	unirseAPartida(cookies[2], t, 1)

	partidaCache := comprobarPartidaEnCurso(t, 1)

	// Simula una partida en la que todos los jugadores tienen un territorio y usuario1 el resto
	for i := logica_juego.Eastern_australia; i <= logica_juego.Alberta; i++ {
		partidaCache.Estado.EstadoMapa[i].NumTropas = 20
		partidaCache.Estado.EstadoMapa[i].Ocupante = "usuario1"
	}

	partidaCache.Estado.EstadoMapa[logica_juego.Eastern_australia].Ocupante = "usuario2"
	partidaCache.Estado.EstadoMapa[logica_juego.Eastern_australia].NumTropas = 1

	partidaCache.Estado.EstadoMapa[logica_juego.Indonesia].Ocupante = "usuario3"
	partidaCache.Estado.EstadoMapa[logica_juego.Indonesia].NumTropas = 1

	globales.CachePartidas.AlmacenarPartida(partidaCache)

	entradasEstado := make([]int, 3)
	for i := 0; i < 3; i++ {
		entradasEstado[i] = len(preguntarEstado(t, cookies[i]).Acciones)
	}

	// Eliminar a usuario2
	atacarYOcupar(t, partidaCache, cookies, logica_juego.Western_australia, logica_juego.Eastern_australia)
	estado := preguntarEstado(t, cookies[1])
	//t.Log(estado.Acciones[len(estado.Acciones)-1])
	t.Log(estado.Acciones)

	encontradaEliminacion := false
	for k, v := range estado.Acciones[len(estado.Acciones)-2].(map[string]interface{}) {
		if k == "IDAccion" {
			if int(v.(float64)) == int(logica_juego.IDAccionJugadorEliminado) {
				encontradaEliminacion = true
				break
			}
		}
	}
	if !encontradaEliminacion {
		t.Fatal("No se ha encontrado una acción de jugador eliminado tras eliminar a usuario2:", estado.Acciones)
	}

	// Eliminar a usuario3
	atacarYOcupar(t, partidaCache, cookies, logica_juego.Western_australia, logica_juego.Indonesia)
	estado = preguntarEstado(t, cookies[2])
	//t.Log(estado.Acciones[len(estado.Acciones)-1])
	t.Log(estado.Acciones)

	encontradaEliminacion = false
	encontradaFinalizacion := false
	for _, accion := range estado.Acciones {
		for k, v := range accion.(map[string]interface{}) {
			if k == "IDAccion" {
				if int(v.(float64)) == int(logica_juego.IDAccionJugadorEliminado) {
					encontradaEliminacion = true
				} else if int(v.(float64)) == int(logica_juego.IDAccionPartidaFinalizada) {
					encontradaFinalizacion = true
				}
			}
		}
	}

	if !encontradaFinalizacion || !encontradaEliminacion {
		t.Fatal("No se ha encontrado un acciones de jugador eliminado y fin de partida tras eliminar a usuario3:", estado.Acciones)
	}

	// Comprobar jugadores restantes por solicitar estado antes de eliminar la partida
	partidaCache = comprobarPartidaEnCurso(t, 1)

	if len(partidaCache.Estado.JugadoresRestantesPorConsultar) != 1 {
		t.Fatal("Número de jugadores por consultar tras haber sido eliminados usuario2 y usuario3 incorrecto:", partidaCache.Estado.JugadoresRestantesPorConsultar)
	} else {
		t.Log("jugadores restantes por preguntar estado antes de cerrar:", partidaCache.Estado.JugadoresRestantesPorConsultar)
	}

	estado = preguntarEstado(t, cookies[0])
	if len(estado.JugadoresRestantesPorConsultar) != 0 {
		t.Fatal("Número de jugadores por consultar tras haber sido eliminados usuario2 y usuario3 y consultar usuario1 incorrecto:", estado.JugadoresRestantesPorConsultar)
	} else {
		t.Log("jugadores restantes por preguntar estado justo antes de cerrar:", estado.JugadoresRestantesPorConsultar)
	}

	_, existePartidaTrasEliminacion := globales.CachePartidas.ObtenerPartida(1)

	if existePartidaTrasEliminacion {
		t.Fatal("La partida aún existe en la cache tras haber sido terminada")
	} else {
		t.Log("Partida eliminada correctamente")
	}

	// Prueba previa a haber permitido que se puedan salir de la partida usuarios eliminados
	/*estado = preguntarEstado(t, cookies[0])
	partidaCache = comprobarPartidaEnCurso(t, 1)
	if len(partidaCache.Estado.JugadoresRestantesPorConsultar) != 1 {
		t.Fatal("Número de jugadores por consultar tras haberlo solicitado usuario2 y usuario1 incorrecto:", partidaCache.Estado.JugadoresRestantesPorConsultar)
	} else {
		t.Log("jugadores restantes por preguntar estado antes de cerrar:", partidaCache.Estado.JugadoresRestantesPorConsultar)
	}
	*/
	/*estado = preguntarEstado(t, cookies[2])
	_, existe := globales.CachePartidas.ObtenerPartida(1)
	if existe {
		t.Fatal("La partida sigue en la cache después de haberse consultado por todos los jugadores tras finalizar")
	}*/

	perfilUsuario1 := obtenerPerfilUsuario(cookies[0], "usuario1", t)
	if perfilUsuario1.PartidasGanadas != 1 {
		t.Fatal("No se ha contabilizado una partida ganada a usuario1:", perfilUsuario1.PartidasGanadas)
	}
	if perfilUsuario1.PartidasTotales != 1 {
		t.Fatal("No se ha contabilizado una partida jugada a usuario1:", perfilUsuario1.PartidasTotales)
	}
	if perfilUsuario1.Puntos != logica_juego.PUNTOS_GANAR {
		t.Fatal("Puntos incorrectos para usuario1", perfilUsuario1.Puntos)
	}

	perfilUsuario2 := obtenerPerfilUsuario(cookies[0], "usuario2", t)
	if perfilUsuario2.PartidasGanadas != 0 {
		t.Fatal("Se ha contabilizado una partida ganada a usuario2:", perfilUsuario2.PartidasGanadas)
	}
	if perfilUsuario2.PartidasTotales != 1 {
		t.Fatal("No se ha contabilizado una partida jugada a usuario2:", perfilUsuario2.PartidasTotales)
	}
	if perfilUsuario2.Puntos != logica_juego.PUNTOS_PERDER {
		t.Fatal("Puntos incorrectos para usuario2", perfilUsuario2.Puntos)
	}

	perfilUsuario3 := obtenerPerfilUsuario(cookies[0], "usuario3", t)
	if perfilUsuario3.PartidasGanadas != 0 {
		t.Fatal("Se ha contabilizado una partida ganada a usuario3:", perfilUsuario3.PartidasGanadas)
	}
	if perfilUsuario3.PartidasTotales != 1 {
		t.Fatal("No se ha contabilizado una partida jugada a usuario3:", perfilUsuario3.PartidasTotales)
	}
	if perfilUsuario3.Puntos != logica_juego.PUNTOS_PERDER {
		t.Fatal("Puntos incorrectos para usuario3", perfilUsuario3.Puntos)
	}
}
