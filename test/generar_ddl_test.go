package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"net/http"
	"testing"
	"time"
)

func TestGenerarPartidaDebug(t *testing.T) {
	//t.Skip("Se ha saltado la generación de ficheros DDL")
	purgarDB()
	// Creamos 6 usuarios, de nombre "jugadorx" y misma contraseña
	jugadores := []string{"jugador1", "jugador2", "jugador3", "jugador4", "jugador5", "jugador6"}
	var cookies []*http.Cookie
	for _, j := range jugadores {
		cookies = append(cookies, crearUsuario(j, t))
	}
	crearPartida(cookies[0], t, true)
	// Unimos todos los jugadores a la partida
	for _, c := range cookies[1:] {
		unirseAPartida(c, t, 1)
	}

	// Comienza la partida, cada uno con 20 tropas
	// Cada uno asigna dichas 20 tropas a uno de sus territorios
	partidaCache := comprobarPartidaEnCurso(t, 1)

	acciones := []interface{}{
		logica_juego.NewAccionRecibirRegion(1, 15, 3, "jugador1"),
		logica_juego.NewAccionCambioFase(1, "jugador1"),
		logica_juego.NewAccionInicioTurno("jugador1", 3, 5, 2),
		logica_juego.NewAccionCambioCartas(1, true, []logica_juego.NumRegion{logica_juego.Afghanistan, logica_juego.Alberta}, false),
		logica_juego.NewAccionReforzar("jugador1", logica_juego.Central_america, 3),
		logica_juego.NewAccionAtaque(logica_juego.Congo, logica_juego.South_africa, 3, 4, 3, "jugador1", "usuario2"),
		logica_juego.NewAccionOcupar(logica_juego.Great_britain, logica_juego.Northern_europe, 2, 7, "jugador1", "jugador2"),
		logica_juego.NewAccionFortificar(logica_juego.China, logica_juego.Alaska, 9, 12, "jugador4"),
		logica_juego.NewAccionObtenerCarta(logica_juego.Carta{IdCarta: 2, Tipo: logica_juego.Artilleria, Region: logica_juego.Northern_europe, EsComodin: false}, "jugador2"),
		logica_juego.NewAccionJugadorEliminado("jugador3", "jugador5", 4),
		logica_juego.NewAccionJugadorExpulsado("jugador5"),
		logica_juego.NewAccionPartidaFinalizada("jugador1"),
	}

	partidaCache.Estado.Acciones = acciones
	globales.CachePartidas.AlmacenarPartida(partidaCache)

	globales.CachePartidas.CanalSerializacion <- partidaCache

	time.Sleep(3 * time.Second)
}

func TestGenerarDDL(t *testing.T) {
	// Comentar la línea 12 para modificar la DB de forma que contenga los datos necesarios para la creación de los ddl
	t.Skip("Se ha saltado la generación de ficheros DDL")

	t.Log("Purgando DB...")
	purgarDB()
	var err error
	// Generar DDL fase de ataque

	// Creamos 6 usuarios, de nombre "jugadorx" y misma contraseña
	jugadores := []string{"jugador1", "jugador2", "jugador3", "jugador4", "jugador5", "jugador6"}
	var cookies []*http.Cookie
	for _, j := range jugadores {
		cookies = append(cookies, crearUsuario(j, t))
	}

	// Creamos una partida pública
	crearPartida(cookies[0], t, true)

	// Unimos todos los jugadores a la partida
	for _, c := range cookies[1:] {
		unirseAPartida(c, t, 1)
	}

	// Comienza la partida, cada uno con 20 tropas
	// Cada uno asigna dichas 20 tropas a uno de sus territorios
	partidaCache := comprobarPartidaEnCurso(t, 1)
	mapa := partidaCache.Estado.EstadoMapa

	// Forzamos a que empiece el jugador1
	saltarTurnos(t, partidaCache, jugadores[0])

	for i, j := range jugadores {
		// Obtengo una región que pertenezca al jugador, y que sea adyacente a una de algún rival
		for r := logica_juego.Eastern_australia; r < logica_juego.Alberta; r++ {
			if mapa[r].Ocupante == j {
				encontrado := false
				adyacentes := logica_juego.Adyacentes(r)
				for _, rr := range adyacentes {
					if mapa[rr].Ocupante != j {
						reforzarTerritorio(t, cookies[i], int(r), 13)
						t.Log("El jugador", j, "ha reforzado el territorio", r, "utilizando 13 tropas")
						encontrado = true
						break
					}
				}

				if encontrado {
					break
				}
			}
		}
		err = saltarFase(cookies[i], t)
		if err != nil {
			t.Fatal("Error al saltar fase:", err)
		}
	}

	// Comienza la fase de ataque
	partidaCache = comprobarPartidaEnCurso(t, 1)
	if partidaCache.Estado.Fase != logica_juego.Ataque {
		t.Fatal("No se ha pasado a la fase de ataque correctamente")
	}
	t.Log("Se ha pasado a la fase de ataque correctamente")
	if partidaCache.Estado.TurnoJugador != 5 {
		t.Fatal("No es el turno del jugador6")
	}
	t.Log("Es el turno del jugador", jugadores[partidaCache.Estado.TurnoJugador])

	time.Sleep(5 * time.Second)
	// Breakpoint aquí para obtener el ddl de la fase de ataque

	// Pasamos a la fase de fortificación
	err = saltarFase(cookies[5], t)
	if err != nil {
		t.Fatal("Error al saltar fase:", err)
	}

	partidaCache = comprobarPartidaEnCurso(t, 1)
	if partidaCache.Estado.Fase != logica_juego.Fortificar {
		t.Fatal("No se ha pasado a la fase de fortificación correctamente")
	}
	t.Log("Se ha pasado a la fase de fortificación correctamente")
	t.Log("Es el turno del jugador", jugadores[partidaCache.Estado.TurnoJugador])

	time.Sleep(5 * time.Second)
	// Breakpoint aquí para obtener el ddl de la fase de fortificación

	// Pasamos a la fase de refuerzo
	err = saltarFase(cookies[5], t)
	if err != nil {
		t.Fatal("Error al saltar fase:", err)
	}
	// Turno del jugador1

	partidaCache = comprobarPartidaEnCurso(t, 1)
	if partidaCache.Estado.Fase != logica_juego.Refuerzo {
		t.Fatal("No se ha pasado a la fase de fortificación correctamente")
	}
	t.Log("Se ha pasado a la fase de fortificación correctamente")
	t.Log("Es el turno del jugador", jugadores[partidaCache.Estado.TurnoJugador])

	// Le damos 6 cartas al jugador, para poder probar cambios
	partidaCache = comprobarPartidaEnCurso(t, 1)
	partidaCache.Estado.EstadosJugadores[jugadores[0]].Cartas = partidaCache.Estado.Cartas[0:6]
	partidaCache.Estado.Cartas = partidaCache.Estado.Cartas[6:]
	globales.CachePartidas.AlmacenarPartida(partidaCache)
	globales.CachePartidas.CanalSerializacion <- partidaCache

	partidaCache = comprobarPartidaEnCurso(t, 1)
	invarianteNumeroDeCartas(partidaCache.Estado, *partidaCache.Estado.EstadosJugadores[jugadores[0]], t)

	// Breakpoint aquí para obtener el ddl de la fase de refuerzo
	time.Sleep(5 * time.Second)
}

func TestComprobarDDL(t *testing.T) {
	// Comentar la siguiente línea para ejecutar el test, que comprueba si la partida 1 cargada en DB
	// corresponde con el ddl previamente generado
	t.Skip("Saltado test para comprobar DDL")
	time.Sleep(5 * time.Second)
	partidaCache := comprobarPartidaEnCurso(t, 1)
	estado := partidaCache.Estado

	// Cambiar en función de la fase inicial de la partida a probar
	fase := logica_juego.Fortificar

	jugadores := estado.Jugadores
	t.Log("Fase:", partidaCache.Estado.Fase, "Turno:", partidaCache.Estado.TurnoJugador)

	if len(jugadores) != 6 {
		t.Fatal("Debería haber 6 jugadores, pero hay", len(jugadores))
	}

	jugadoresEsperados := []string{"jugador1", "jugador2", "jugador3", "jugador4", "jugador5", "jugador6"}
	for i := 0; i < 6; i++ {
		if jugadores[i] != jugadoresEsperados[i] {
			t.Fatal("El jugador", jugadores[i], "no coincide con", jugadoresEsperados[i])
		}
	}

	t.Log("Los jugadores coinciden con los esperados:", jugadores)

	switch fase {
	case logica_juego.Refuerzo:
		if estado.Fase != logica_juego.Refuerzo {
			t.Fatal("La fase debería ser refuerzo, pero es:", estado.Fase)
		}

		if estado.TurnoJugador != 0 {
			t.Fatal("El turno debería ser 0, pero es:", estado.TurnoJugador)
		}

		if len(estado.EstadosJugadores[jugadores[0]].Cartas) != 6 {
			t.Fatal("El jugador1 debería tener 6 cartas")
		}
		t.Log("OK, partida cargada en fase Refuerzo, turno del jugador1, con 6 cartas")
	case logica_juego.Ataque:
		if estado.Fase != logica_juego.Ataque {
			t.Fatal("La fase debería ser ataque, pero es:", estado.Fase)
		}

		if estado.TurnoJugador != 5 {
			t.Fatal("El turno debería ser 5, pero es:", estado.TurnoJugador)
		}

		t.Log("OK, partida cargada en fase Ataque, turno del jugador6")
	case logica_juego.Fortificar:
		if estado.Fase != logica_juego.Fortificar {
			t.Fatal("La fase debería ser fortificar, pero es:", estado.Fase)
		}

		if estado.TurnoJugador != 5 {
			t.Fatal("El turno debería ser 5, pero es:", estado.TurnoJugador)
		}

		t.Log("OK, partida cargada en fase Fortificar, turno del jugador6")
	}
}
