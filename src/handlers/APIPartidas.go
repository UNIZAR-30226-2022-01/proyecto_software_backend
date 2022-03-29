package handlers

import (
	"backend/dao"
	"backend/globales"
	"backend/middleware"
	"backend/vo"
	"database/sql"
	"encoding/json"
	"errors"
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
	if maxJugadores < 2 || maxJugadores > 6 {
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

	devolverExito(writer)
}

// UnirseAPartida permite al usuario unirse a una partida en caso de que no esté en otra,
// no esté completa la partida, sea pública, o tenga su contraseña si es privada.
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

	// TODO: Probar
	// Ya se puede empezar la partida
	if (len(jugadores) + 1) == maxJugadores {
		partida, err = dao.ObtenerPartida(globales.Db, idPartida)

		// Se añade al usuario e inicia, creando su estado
		jugadores = append(jugadores, usuario)
		partida.IniciarPartida(jugadores)

		// Se añade al almacén
		globales.CachePartidas.AlmacenarPartida(partida)

		// Y se encola un trabajo de serialización de su estado
		globales.CachePartidas.CanalSerializacion <- partida
	}

	devolverExito(writer)
}

// AbandonarLobby deja la partida en la que el usuario esté participando. Responde con status 200 si ha habido éxito,
// o status 500 si ha habido un error junto a su motivo en el cuerpo.
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
		devolverExito(writer)
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
	}
}

// TODO: documentar y probar
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
}

// TODO: documentar y probar
func ReforzarTerritorio(writer http.ResponseWriter, request *http.Request) {
	idTerritorio, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		devolverError(writer, errors.New("Se ha introducido un identificador de territorio inválido: "+chi.URLParam(request, "id")))
	}

	numTropas, err := strconv.Atoi(chi.URLParam(request, "numTropas"))
	if err != nil {
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
	}
}
