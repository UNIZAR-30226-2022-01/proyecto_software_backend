// Package handlers define las funciones de tratamiento de cada ruta a la API
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// CrearPartida crea una nueva partida, para la que se definirá el número máximo de jugadores,
// si es pública o privada, y la contraseña en caso de que fuera necesario
// Parámetros del formulario recibido:
//	El campo "maxJugadores" indica el número máximo de jugadores, "tipo"	indica si la partida es pública o privada.
//	Si la cadena "tipo" equivale a "Publica", la partida será pública, en cualquier otro caso será privada.
//  El campo "password" define la contraseña necesaria para el acceso a una partida privada.
//
// Si se produjera un error durante el procesado, se devolverá código 500
// En cualquier otro caso, se devolverá código 200
//
// Ruta: /api/crearPartida
// Tipo: POST
func CrearPartida(writer http.ResponseWriter, request *http.Request) {
	password := request.FormValue("password")
	maxJugadores, err := strconv.Atoi(request.FormValue("maxJugadores"))
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	tipoPartida := request.FormValue("tipo")
	esPublica := tipoPartida == "Publica"

	if err != nil {
		devolverError(writer, errors.New("Se ha introducido un valor no numérico en el número de jugadores."))
		return
	}
	if maxJugadores < 3 || maxJugadores > 6 {
		devolverError(writer, errors.New("El número de jugadores debe ser un valor numérico entre 2 y 6."))
		return
	}

	hash := ""
	if !esPublica {
		hash, err = hashPassword(password)
		if err != nil {
			devolverError(writer, errors.New("Se ha producido un error al procesar los datos."))
			return
		}
		hash = hash
	}

	usuario := vo.Usuario{"", nombreUsuario, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	//partida = vo.Partida{0, esPublica, partida.PasswordHash, false, maxJugadores, nil, []vo.Mensaje{}, vo.EstadoPartida{}}

	enPartida, err := dao.UsuarioEnPartida(globales.Db, &usuario)
	if enPartida {
		devolverError(writer, errors.New("Ya estás participando en otra partida."))
		return
	}

	partida := *vo.CrearPartida(esPublica, hash, maxJugadores)

	err = dao.CrearPartida(globales.Db, &usuario, &partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	globales.CachePartidas.AlmacenarPartida(partida)

	escribirHeaderExito(writer)
}

// UnirseAPartida permite al usuario unirse a una partida en caso de que no esté en otra,
// no esté completa la partida, sea pública, o tenga su contraseña si es privada.
// Si se produciera algún error, devuelve código 500, en caso contrario 200
// Los campos del formulario son "password" e "idPartida"
//
// Ruta: /api/unirseAPartida
// Tipo: POST
func UnirseAPartida(writer http.ResponseWriter, request *http.Request) {
	password := request.FormValue("password")
	idPartida, err := strconv.Atoi(request.FormValue("idPartida"))
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	usuario := vo.Usuario{NombreUsuario: nombreUsuario}
	partida := vo.Partida{IdPartida: idPartida}
	jugadores, maxJugadores, err := dao.ConsultarJugadoresPartida(globales.Db, &partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Comprobamos que la partida no esté completa (Puede haber intentado entrar justo antes de haber empezado la partida)
	if len(jugadores) == maxJugadores {
		devolverError(writer, errors.New("No hay hueco en la partida."))
		return
	}

	// Comprobamos que el usuario no esté participando en otra partida
	enPartida, err := dao.UsuarioEnPartida(globales.Db, &usuario)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	if enPartida {
		devolverError(writer, errors.New("El usuario ya está en otra partida"))
		return
	}

	publica, passwordHash, err := dao.ConsultarAcceso(globales.Db, &partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	if !publica {
		// Comprobamos que la contraseña sea correcta
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		if err != nil {
			devolverError(writer, errors.New("La contraseña no es correcta."))
			return
		}
	}

	// Else -> no está completa, el usuario no está en otra partida y la partida es pública o la contraseña es correcta
	err = dao.UnirseAPartida(globales.Db, &usuario, &partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Ya se puede empezar la partida
	if (len(jugadores) + 1) == maxJugadores {
		partida, err = dao.ObtenerPartida(globales.Db, idPartida)
		if err != nil {
			devolverErrorSQL(writer)
		}

		err = dao.EmpezarPartida(globales.Db, idPartida)
		if err != nil {
			devolverErrorSQL(writer)
		}

		// Se añade al usuario e inicia, creando su estado
		jugadores = append(jugadores, usuario)

		nombresJugadores := []string{}
		for _, j := range jugadores {
			nombresJugadores = append(nombresJugadores, j.NombreUsuario)
		}

		partida.IniciarPartida(nombresJugadores)

		// Se añade al almacén
		globales.CachePartidas.AlmacenarPartida(partida)

		// Y se encola un trabajo de serialización de su estado
		globales.CachePartidas.CanalSerializacion <- partida
	}

	escribirHeaderExito(writer)
}

// ObtenerEstadoLobby devuelve el estado del lobby de una partida identificada por su id
// Devuelve si es pública o no, si está o no en curso, el número máximo de jugadores y
// los jugadores que se encuentran en el lobby
// Devuelve código de error 500 en caso de error, código 200 en cualquier otro caso
// El JSON devuelto tiene el siguiente formato
// {
//  "EnCurso":bool,
// 	"EsPublico":bool,
//  "Jugadores":int,
//  "MaxJugadores":int,
//  "NombresJugadores": [string, string, ...]
// }
//
// Ruta: /api/obtenerEstadoLobby/{id}
// Tipo: GET
func ObtenerEstadoLobby(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	partida, err := dao.ObtenerIDPartida(globales.Db, nombreUsuario)
	if err != nil {
		devolverError(writer, err)
		return
	}

	estadoLobby, err := dao.ObtenerEstadoLobby(globales.Db, partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(estadoLobby)
	escribirHeaderExito(writer)
}

// AbandonarLobby deja la partida en la que el usuario esté participando. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
//
// Ruta: /api/abandonarLobby
// Tipo: POST
func AbandonarLobby(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)

	err := dao.AbandonarLobby(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})
	// AbandonarLobby ya da el error formateado
	if err != nil {
		devolverError(writer, err)
	} else {
		escribirHeaderExito(writer)
	}
}

// ObtenerPartidas devuelve un listado de partidas codificado en JSON, con el siguiente orden:
//	1- partidas privadas, de más a menos amigos presentes
//	2- partidas públicas, de más a menos amigos presentes
//	3- partidas públicas sin amigos: de más a menos jugadores
//	4- partidas privadas sin amigos: de más a menos jugadores
//
// El contenido JSON tendrá la siguiente forma:
//
// [												// Lista de elementos de partida (Nota: Será nula si no existen partidas)
//
//													// Elemento de partida
//  {
//    "IdPartida": int,							 	// ID de partida (Nota: Se deberá usar como referencia en el resto de peticiones relativas a ella)
//    "EsPublica": bool,							// Flag de partida pública o privada
//    "NumeroJugadores": int,						// Número de jugadores presentes en el lobby
//    "MaxNumeroJugadores": int,					// Número de jugadores máximo establecido
//    "AmigosPresentes": [ string, string, ...], 	// Lista de nombres de amigos presentes en el lobby (Nota: Será nulo si "NumAmigosPresentes" tiene como valor 0)
//    "NumAmigosPresentes": int						// Número de amigos presentes en el lobby
//  },
//
//  {...}
// ]
//
//
//
// Ejemplo:
// [
//  {
//    "IdPartida": 1,
//    "EsPublica": false,
//    "NumeroJugadores": 4,
//    "MaxNumeroJugadores": 6,
//    "AmigosPresentes": [
//      "amigo1",
//      "amigo2"
//    ],
//    "NumAmigosPresentes": 2
//  },
//  {
//    "IdPartida": 2,
//    "EsPublica": false,
//    "NumeroJugadores": 4,
//    "MaxNumeroJugadores": 6,
//    "AmigosPresentes": [
//      "amigo3"
//    ],
//    "NumAmigosPresentes": 1
//  },
//  {
//    "IdPartida": 3,
//    "EsPublica": true,
//    "NumeroJugadores": 3,
//    "MaxNumeroJugadores": 6,
//    "AmigosPresentes": null,
//    "NumAmigosPresentes": 0
//  }
//]
//
// Si ocurre algún error durante el procesamiento, se devolverá un status 500.
//
// Ruta: /api/obtenerPartidas
// Tipo: GET
func ObtenerPartidas(writer http.ResponseWriter, request *http.Request) {
	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	amigos, err := dao.ObtenerAmigos(globales.Db, &usuario)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	partidas, err := dao.ObtenerPartidasNoEnCurso(globales.Db)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	partidasPrivadas, partidasPublicas := dividirPartidasPrivadasYPublicas(partidas)

	// Ordena partidas privadas de más a menos amigos
	ordenarPorNumeroAmigos(partidasPrivadas, amigos)
	// Ordena partidas públicas de más a menos amigos
	ordenarPorNumeroAmigos(partidasPublicas, amigos)
	// Extrae las partidas privadas sin amigos del usuario del slice y deja las partidas privadas con amigos
	partidasPrivadasConAmigos, partidasPrivadasSinAmigos := dividirPartidasPorAmigos(partidasPrivadas, amigos)

	// Ordena partidas privadas sin amigos de más a menos jugadores
	ordenarPorNumeroJugadores(writer, partidasPrivadasSinAmigos)
	// Extrae las partidas públicas sin amigos del usuario del slice y deja las partidas públicas con amigos
	partidasPublicasConAmigos, partidasPublicasSinAmigos := dividirPartidasPorAmigos(partidasPublicas, amigos)

	// Ordena partidas públicas sin amigos de más a menos jugadores
	ordenarPorNumeroJugadores(writer, partidasPublicasSinAmigos)

	// Junta todos los slices, en orden
	var partidasOrdenadas []vo.Partida
	partidasOrdenadas = append(partidasOrdenadas, partidasPrivadasConAmigos...)
	partidasOrdenadas = append(partidasOrdenadas, partidasPublicasConAmigos...)
	partidasOrdenadas = append(partidasOrdenadas, partidasPublicasSinAmigos...)
	partidasOrdenadas = append(partidasOrdenadas, partidasPrivadasSinAmigos...)

	if err != nil {
		devolverErrorSQL(writer)
	} else {
		var elementos []vo.ElementoListaPartidas

		for _, p := range partidasOrdenadas {
			elementos = append(elementos, transformarAElementoListaPartidas(&p, amigos))
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(elementos)
		escribirHeaderExito(writer)
	}
}

// ObtenerEstadoPartida devuelve la lista de acciones transcurridas desde la última consulta del usuario hasta
// el momento, que deberán ser procesadas en orden.
// El formato es una lista de acciones, codificada en JSON de la siguiente forma:
// [{acción}, {acción}]
//
// Donde cada acción es una acción específica a distinguir según el primer campo común a todas, "IDAccion", para su interpretación.
//
// Ejemplo:
// [
//   {
//      "IDAccion":0,
//      "Region":0,
//      "TropasRestantes":19,
//      "TerritoriosRestantes":41,
//      "Jugador":"usuario4"
//   },
//   {
//      "IDAccion":1,
//      "Jugador":"usuario5"
//   }
//]
//
// La lista de acciones y su formato en JSON están disponibles en el módulo de logica_juego, en acciones.go
//
// Ruta: /api/obtenerEstadoPartida
// Tipo: GET
func ObtenerEstadoPartida(writer http.ResponseWriter, request *http.Request) {
	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	idPartida, err := dao.PartidaUsuario(globales.Db, &usuario)
	if err == sql.ErrNoRows {
		devolverError(writer, errors.New("No estás participando en ninguna partida."))
		return
	} else if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Se obtiene una copia de la partida
	partida, existe := globales.CachePartidas.ObtenerPartida(idPartida)
	if !existe {
		devolverErrorSQL(writer)
		return
	}

	// Se obtiene acciones de [UltimoIndiceLeido+1...)
	acciones := partida.Estado.Acciones[partida.Estado.EstadosJugadores[usuario.NombreUsuario].UltimoIndiceLeido+1:]

	// Se marca que el usuario ha leído hasta el último índice
	partida.Estado.EstadosJugadores[usuario.NombreUsuario].UltimoIndiceLeido = len(partida.Estado.Acciones) - 1

	// Y si ha terminado, o el jugador ha perdido
	if partida.Estado.Terminada || partida.Estado.ContarTerritoriosOcupados(usuario.NombreUsuario) == 0 {
		for i, jugador := range partida.Estado.JugadoresRestantesPorConsultar {
			// Si aún no había comprobado el estado hasta ahora
			if jugador == usuario.NombreUsuario {
				err := terminarPartida(usuario, &partida, i)
				if err != nil {
					devolverErrorSQL(writer)
					return // No se procesará el potencial fin de partida o modifica en la BD/cache si hay un error
				}
			}
		}
	}

	// Si ha terminado y no queda ningún usuario más por consultar su estado, se elimina
	if partida.Estado.Terminada && (len(partida.Estado.JugadoresRestantesPorConsultar) == 0) {
		// De la DB
		err := dao.BorrarPartida(globales.Db, &vo.Partida{IdPartida: idPartida})
		if err != nil {
			log.Println("Error al borrar una partida de la base de datos, borrando del cache:", err)
		}

		// De la cache
		globales.CachePartidas.EliminarPartida(vo.Partida{IdPartida: idPartida})
	} else {
		// Se sobreescribe en el almacén
		globales.CachePartidas.AlmacenarPartida(partida)

		// Y se encola un trabajo de serialización de su estado
		globales.CachePartidas.CanalSerializacion <- partida
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(acciones)
	escribirHeaderExito(writer)
}

// ObtenerEstadoPartidaCompleto devuelve la lista de todas las acciones transcurridas desde el inicio de la partida
// hasta el momento, que deberán ser procesadas en orden. Las siguientes llamadas a ObtenerEstadoPartida consultarán las
// acciones desde dicho momento, pudiendo ser por tanto un sustituto a ObtenerEstadoPartida.
// El formato es una lista de acciones, codificada en JSON de la siguiente forma:
// [{acción}, {acción}]
//
// Donde cada acción es una acción específica a distinguir según el primer campo común a todas, "IDAccion", para su interpretación.
//
// Ejemplo:
// [
//   {
//      "IDAccion":0,
//      "Region":0,
//      "TropasRestantes":19,
//      "TerritoriosRestantes":41,
//      "Jugador":"usuario4"
//   },
//   {
//      "IDAccion":1,
//      "Jugador":"usuario5"
//   }
//]
//
// La lista de acciones y su formato en JSON están disponibles en el módulo de logica_juego, en acciones.go
//
// Ruta: /api/obtenerEstadoPartida
// Tipo: GET
func ObtenerEstadoPartidaCompleto(writer http.ResponseWriter, request *http.Request) {
	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	idPartida, err := dao.PartidaUsuario(globales.Db, &usuario)
	if err == sql.ErrNoRows {
		devolverError(writer, errors.New("No estás participando en ninguna partida."))
		return
	} else if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Se obtiene una copia de la partida
	partida, existe := globales.CachePartidas.ObtenerPartida(idPartida)
	if !existe {
		devolverErrorSQL(writer)
		return
	}

	// Se obtiene acciones de [UltimoIndiceLeido+1...)
	acciones := partida.Estado.Acciones

	// Se marca que el usuario ha leído hasta el último índice
	partida.Estado.EstadosJugadores[usuario.NombreUsuario].UltimoIndiceLeido = len(partida.Estado.Acciones) - 1

	// Y si ha terminado, o el jugador ha perdido
	if partida.Estado.Terminada || partida.Estado.ContarTerritoriosOcupados(usuario.NombreUsuario) == 0 {
		for i, jugador := range partida.Estado.JugadoresRestantesPorConsultar {
			// Si aún no había comprobado el estado hasta ahora
			if jugador == usuario.NombreUsuario {
				err := terminarPartida(usuario, &partida, i)
				if err != nil {
					devolverErrorSQL(writer)
					return // No se procesará el potencial fin de partida o modifica en la BD/cache si hay un error
				}
			}
		}
	}

	// Si ha terminado y no queda ningún usuario más por consultar su estado, se elimina
	if partida.Estado.Terminada && (len(partida.Estado.JugadoresRestantesPorConsultar) == 0) {
		// De la DB
		err := dao.BorrarPartida(globales.Db, &vo.Partida{IdPartida: idPartida})
		if err != nil {
			log.Println("Error al borrar una partida de la base de datos, borrando del cache:", err)
		}

		// De la cache
		globales.CachePartidas.EliminarPartida(vo.Partida{IdPartida: idPartida})
	} else {
		// Se sobreescribe en el almacén
		globales.CachePartidas.AlmacenarPartida(partida)

		// Y se encola un trabajo de serialización de su estado
		globales.CachePartidas.CanalSerializacion <- partida
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(acciones)
	escribirHeaderExito(writer)
}

// JugandoEnPartida devuelve verdad si el jugador está participando en una partida, o falso en caso contrario.
// El formato es un booleano codificado en JSON de la siguiente forma:
// true
// false
//
// Ruta: /api/jugandoEnPartida
// Tipo: GET
func JugandoEnPartida(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)

	esta, err := dao.UsuarioEnPartida(globales.Db, &vo.Usuario{NombreUsuario: nombreUsuario})

	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(esta)
	escribirHeaderExito(writer)
}

// Trata el abandono de una partida por parte de un jugador dado, dejando de participar en ella,
// otorgando los puntos según haya ganado o perdido, contabilizando que el usuario ha participado en una
// partida y que ya ha consultado el estado final en el estado de la partida
func terminarPartida(usuario vo.Usuario, partida *vo.Partida, i int) error {
	err := dao.AbandonarPartida(globales.Db, &usuario)

	if err != nil && err != sql.ErrNoRows {
		// Es un error, pero no derivado de un reintento de consultar acciones e intentar abandonar
		// tras un error previo (se ha ejecutado el abandono pero ha fallado otorgar puntos)
		log.Println("Error al abandonar una partida, forzando un reintento más tarde:", err)
		return err
	} else {
		// Otorga al jugador puntos dependiendo de cómo haya quedado en la partida
		if partida.Estado.ContarTerritoriosOcupados(usuario.NombreUsuario) == 0 { // Ha perdido
			err = dao.OtorgarPuntos(globales.Db, &usuario, logica_juego.PUNTOS_PERDER, false)
			if err != nil {
				// Fuerza a que el jugador consulte el estado más tarde para poder salir, al no registrarlo
				log.Println("Error al otorgar puntos a", usuario, ":", err)
				return err
			}
			err = dao.ContabilizarPartida(globales.Db, &usuario)
		} else { // Ha ganado
			err = dao.OtorgarPuntos(globales.Db, &usuario, logica_juego.PUNTOS_GANAR, true)
			if err != nil {
				// Fuerza a que el jugador consulte el estado más tarde para poder salir, al no registrarlo
				log.Println("Error al otorgar puntos a", usuario, ":", err)
				return err
			}
			err = dao.ContabilizarPartidaGanada(globales.Db, &usuario)
		}

		if err != nil {
			// Fuerza a que el jugador consulte el estado más tarde para poder salir, al no registrarlo
			log.Println("Error al contabilizar partida jugada/ganada a", usuario, ":", err)
			return err
		} else {
			// Lo registra
			partida.Estado.JugadoresRestantesPorConsultar = append(partida.Estado.JugadoresRestantesPorConsultar[:i], partida.Estado.JugadoresRestantesPorConsultar[i+1:]...)
		}
	}

	return err
}
