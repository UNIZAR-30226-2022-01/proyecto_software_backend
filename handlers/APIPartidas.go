// Package handlers define las funciones de tratamiento de cada ruta a la API
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// CrearPartida crea una nueva partida, para la que se definirá el número máximo de jugadores,
// si es pública o privada, y la contraseña en caso de que fuera necesario
// Parámetros del formulario recibido:
//	"maxJugadores" indica el número máximo de jugadores
//	"tipo"	indica si la partida es pública o privada
//		si tipo== "Publica", será publica, en cualquier otro caso será privada
//  "password" define la contraseña necesaria para el acceso a una partida privada
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
// [
//  "EnCurso":bool
// 	"EsPublico":bool
//  "Jugadores":int
//  "MaxJugadores":int
//  "NombresJugadores": [string, string, ...]
// ]
//
// Ruta: /api/obtenerEstadoLobby/{id}
// Tipo: GET
func ObtenerEstadoLobby(writer http.ResponseWriter, request *http.Request) {
	partida, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		devolverError(writer, errors.New("El id de la partida debe ser un número entero"))
	}

	estadoLobby, err := dao.ObtenerEstadoLobby(globales.Db, partida)
	if err != nil {
		devolverErrorSQL(writer)
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
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write([]byte(err.Error()))
		if err != nil {
			log.Println("Error al escribir respuesta en:", err)
		}
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
	}

	partidas, err := dao.ObtenerPartidasNoEnCurso(globales.Db)
	if err != nil {
		devolverErrorSQL(writer)
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
	} else if err != nil {
		devolverErrorSQL(writer)
	}

	// Se obtiene una copia de la partida
	partida, existe := globales.CachePartidas.ObtenerPartida(idPartida)
	if !existe {
		devolverErrorSQL(writer)
	}

	// Indexado de slices en go: [x:y)
	// Se obtiene acciones de [UltimoIndiceLeido+1...)
	acciones := partida.Estado.Acciones[partida.Estado.EstadosJugadores[usuario.NombreUsuario].UltimoIndiceLeido+1:]

	// Se marca que el usuario ha leído hasta el último índice
	partida.Estado.EstadosJugadores[usuario.NombreUsuario].UltimoIndiceLeido = len(partida.Estado.Acciones) - 1

	// Se sobreescribe en el almacén
	globales.CachePartidas.AlmacenarPartida(partida)

	// Y se encola un trabajo de serialización de su estado
	globales.CachePartidas.CanalSerializacion <- partida

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(acciones)
	escribirHeaderExito(writer)
}

// ReforzarTerritorio refuerza un territorio con su identificador numérico "id" con un valor de tropas numérico
// codificado en "numTropas", ambos parámetros de la URL.
//
// En caso de éxito, se devolverá un código HTTP 200 y aparecerá una nueva acción "AccionReforzar" en la siguiente
// consulta al estado indicando los detalles de la acción realizada.
//
// En caso de error (número de tropas incorrecto, el turno del jugador es incorrecto, etc.) se devolverá un código HTTP
// 500 junto al mensaje de error en el cuerpo.
//
// Ruta: /reforzarTerritorio/{id}/{numTropas}
// Tipo: POST
func ReforzarTerritorio(writer http.ResponseWriter, request *http.Request) {
	idTerritorio, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		devolverError(writer, errors.New("Se ha introducido un identificador de territorio inválido: "+chi.URLParam(request, "id")))
	}

	numTropas, err := strconv.Atoi(chi.URLParam(request, "numTropas"))
	if err != nil || numTropas <= 0 {
		devolverError(writer, errors.New("Se ha introducido un número de tropas inválido: "+chi.URLParam(request, "numTropas")))
	}

	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	idPartida, err := dao.PartidaUsuario(globales.Db, &usuario)
	if err == sql.ErrNoRows {
		devolverError(writer, errors.New("No estás participando en ninguna partida."))
	} else if err != nil {
		devolverErrorSQL(writer)
	}

	partida, _ := globales.CachePartidas.ObtenerPartida(idPartida)

	err = partida.Estado.ReforzarTerritorio(idTerritorio, numTropas, usuario.NombreUsuario)
	if err != nil {
		devolverError(writer, err)
	} else {
		// Se sobreescribe en el almacén
		globales.CachePartidas.AlmacenarPartida(partida)

		// Y se encola un trabajo de serialización de su estado
		globales.CachePartidas.CanalSerializacion <- partida
		escribirHeaderExito(writer)
	}
}

// CambiarCartas permite a un jugador cambiar un conjunto de 3 cartas por tropas. El número de tropas recibidas
// dependerá del número de cambios totales realizados:
// 		- En el primer cambio se recibirán 4 cartas
//		- Por cada cambio, se recibirán 2 cartas más que en el anterior
//		- En el sexto cambio se recibirán 15 cartas
// 		- A partir del sexto cambio, se recibirán 5 cartas más que en el cambio anterior
//
// Los cambios válidos son los siguientes:
//		- 3 cartas del mismo tipo
//		- 2 cartas del mismo tipo más un comodín
//		- 3 cartas, una de cada tipo
//
// Si el jugador cambia una carta en la que aparece un territorio ocupado por él, se añadirán dos tropas a ese territorio.
// Ruta: /api/cambiarCartas/{carta1}/{carta2}/{carta3}/
// Tipo: GET
func CambiarCartas(writer http.ResponseWriter, request *http.Request) {
	idCarta1, err1 := strconv.Atoi(chi.URLParam(request, "carta1"))
	idCarta2, err2 := strconv.Atoi(chi.URLParam(request, "carta2"))
	idCarta3, err3 := strconv.Atoi(chi.URLParam(request, "carta3"))

	if err1 != nil || err2 != nil || err3 != nil {
		devolverError(writer, errors.New("Los identificadores de las cartas deben ser números naturales"))
		return
	}
	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	idPartida, err := dao.PartidaUsuario(globales.Db, &usuario)
	if err == sql.ErrNoRows {
		devolverError(writer, errors.New("No estás participando en ninguna partida."))
		return
	} else if err != nil {
		devolverErrorSQL(writer)
		return
	}

	partida, _ := globales.CachePartidas.ObtenerPartida(idPartida)

	err = partida.Estado.CambiarCartas(usuario.NombreUsuario, idCarta1, idCarta2, idCarta3)
	if err != nil {
		devolverError(writer, err)
	} else {
		// Se sobreescribe en el almacén
		globales.CachePartidas.AlmacenarPartida(partida)

		// Y se encola un trabajo de serialización de su estado
		globales.CachePartidas.CanalSerializacion <- partida
		escribirHeaderExito(writer)
	}
}

// ConsultarCartas permite al usuario consultar las cartas que tiene en la mano mientras juega una partida
// Un usuario podrá consultar únicamente sus propias cartas.
// El JSON enviado como respuesta tendrá el siguiente formato:
// TODO formato JSON
// Ruta: /api/consultarCartas
// Tipo: GET
func ConsultarCartas(writer http.ResponseWriter, request *http.Request) {
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	usuario := vo.Usuario{NombreUsuario: nombreUsuario}

	idPartida, err := dao.PartidaUsuario(globales.Db, &usuario)
	if err == sql.ErrNoRows {
		devolverError(writer, errors.New("No estás participando en ninguna partida."))
		return
	} else if err != nil {
		devolverErrorSQL(writer)
		return
	}

	partida, _ := globales.CachePartidas.ObtenerPartida(idPartida)
	cartas := partida.Estado.ConsultarCartas(nombreUsuario)

	// Enviamos como respuesta un JSON que contenga las cartas
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(cartas)
	escribirHeaderExito(writer)
}
