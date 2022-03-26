package handlers

import (
	"backend/dao"
	"backend/globales"
	"backend/middleware"
	"backend/vo"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

// CrearPartida crea una nueva partida, para la que se definirá el número máximo de jugadores,
// si es pública o privada, y la contraseña en caso de que fuera necesario
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

	var partida vo.Partida
	hash := ""
	if !esPublica {
		hash, err = hashPassword(password)
		if err != nil {
			devolverError(writer, errors.New("Se ha producido un error al procesar los datos."))
			return
		}
		partida.PasswordHash = hash
	}

	usuario := vo.Usuario{"", nombreUsuario, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	partida = vo.Partida{0, esPublica, partida.PasswordHash, false, maxJugadores, nil, []vo.Mensaje{}, vo.EstadoPartida{}}
	partida.CrearEstadoPartida()
	partida.InicializarAcciones()

	partida.Jugadores = make([]vo.Usuario, 6)
	partida.Jugadores = append(partida.Jugadores, usuario)

	enPartida, err := dao.UsuarioEnPartida(globales.Db, &usuario)
	if enPartida {
		devolverError(writer, errors.New("Ya estás participando en otra partida."))
		return
	}

	err = dao.CrearPartida(globales.Db, &usuario, &partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

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
	jugadores, maxJugadores, err := dao.ConsultarNumeroJugadores(globales.Db, &partida)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Comprobamos que la partida no esté completa
	if jugadores == maxJugadores {
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

	partidas, err := dao.ObtenerPartidas(globales.Db)
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

func ordenarPorNumeroJugadores(writer http.ResponseWriter, partidasPrivadasSinAmigos []vo.Partida) {
	sort.SliceStable(partidasPrivadasSinAmigos, func(i, j int) bool {
		// Orden: > a <
		numeroJugadoresI, _, err1 := dao.ConsultarNumeroJugadores(globales.Db, &partidasPrivadasSinAmigos[i])
		numeroJugadoresJ, _, err2 := dao.ConsultarNumeroJugadores(globales.Db, &partidasPrivadasSinAmigos[i])

		if err1 != nil || err2 != nil {
			devolverErrorSQL(writer)
		}

		return numeroJugadoresI > numeroJugadoresJ
	})
}

func dividirPartidasPorAmigos(partidasPrivadas []vo.Partida, amigos []vo.Usuario) ([]vo.Partida, []vo.Partida) {
	var partidasPrivadasConAmigos []vo.Partida
	var partidasPrivadasSinAmigos []vo.Partida
	for _, partida := range partidasPrivadas {
		// Se ha llegado al punto en el slice a partir del cual no hay amigos
		if vo.ContarAmigos(amigos, partida) == 0 {
			partidasPrivadasSinAmigos = append(partidasPrivadasSinAmigos, partida)
		} else {
			partidasPrivadasConAmigos = append(partidasPrivadasConAmigos, partida)
		}
	}
	return partidasPrivadasConAmigos, partidasPrivadasSinAmigos
}

func ordenarPorNumeroAmigos(partidasPrivadas []vo.Partida, amigos []vo.Usuario) {
	sort.SliceStable(partidasPrivadas, func(i, j int) bool {
		// Orden: > a <
		return vo.ContarAmigos(amigos, partidasPrivadas[i]) > vo.ContarAmigos(amigos, partidasPrivadas[j])
	})
}

func dividirPartidasPrivadasYPublicas(partidas []vo.Partida) ([]vo.Partida, []vo.Partida) {
	// Extrae las partidas privadas del slice y deja las partidas públicas
	var partidasPrivadas []vo.Partida
	var partidasPublicas []vo.Partida
	for _, partida := range partidas {
		if !partida.EsPublica {
			partidasPrivadas = append(partidasPrivadas, partida)
		} else {
			partidasPublicas = append(partidasPublicas, partida)
		}
	}
	return partidasPrivadas, partidasPublicas
}

// transformarAElementoListaPartidas convierte una partida en un elemento de lista de partidas,
// dada una lista de amigos de un usuario. Se asume que la partida existe en la DB.
// No puede localizarse en el módulo VO porque causaría una dependencia cíclica con DAO
func transformarAElementoListaPartidas(p *vo.Partida, amigos []vo.Usuario) vo.ElementoListaPartidas {
	listaAmigos := obtenerAmigos(amigos, p)
	numeroJugadores, _, err := dao.ConsultarNumeroJugadores(globales.Db, p)
	if err != nil {
		numeroJugadores = 0
	}

	return vo.ElementoListaPartidas{
		IdPartida:          p.IdPartida,
		EsPublica:          p.EsPublica,
		NumeroJugadores:    numeroJugadores,
		MaxNumeroJugadores: p.MaxNumeroJugadores,
		AmigosPresentes:    listaAmigos,
		NumAmigosPresentes: len(listaAmigos),
	}
}

// obtenerAmigos obtiene una lista de nombres amigos presentes en
// una partida, dada una lista previa
func obtenerAmigos(amigos []vo.Usuario, partida *vo.Partida) (listaFiltrada []string) {
	for _, amigo := range amigos {
		// Como máximo hay 6 jugadores en la partida, así que
		// la complejidad la dicta el número de amigos del usuario
		for _, jugador := range partida.Jugadores {
			if amigo.NombreUsuario == jugador.NombreUsuario {
				listaFiltrada = append(listaFiltrada, amigo.NombreUsuario)
			}
		}
	}

	return listaFiltrada
}
