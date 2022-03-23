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
		devolverError(writer, "CrearPartida", err)
		return
	}
	if maxJugadores < 2 || maxJugadores > 6 {
		devolverError(writer, "CrearPartida", errors.New("el número de jugadores debe estar entre 2 y 6"))
		return
	}

	var partida vo.Partida
	hash := ""
	if !esPublica {
		hash, err = hashPassword(password)
		partida.PasswordHash = hash
	}
	log.Println("Partida publica", esPublica, "hash:", partida.PasswordHash)
	if err != nil {
		devolverError(writer, "CrearPartida", err)
		return
	}

	usuario := vo.Usuario{"", nombreUsuario, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	partida = vo.Partida{0, esPublica, partida.PasswordHash, false, maxJugadores, nil, []vo.Mensaje{}, vo.EstadoPartida{}}
	partida.Jugadores = make([]vo.Usuario, 6)
	partida.Jugadores = append(partida.Jugadores, usuario)

	enPartida, err := dao.UsuarioEnPartida(globales.Db, &usuario)
	if enPartida {
		devolverError(writer, "Crear Partida", errors.New("El usuario ya está participando en otra partida"))
		return
	}

	err = dao.CrearPartida(globales.Db, &usuario, &partida)
	if err != nil {
		devolverError(writer, "Crear Partida", err)
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
		devolverError(writer, "Unirse a Partida", err)
		return
	}

	usuario := vo.Usuario{NombreUsuario: nombreUsuario}
	partida := vo.Partida{IdPartida: idPartida}
	jugadores, maxJugadores, err := dao.ConsultarNumeroJugadores(globales.Db, &partida)
	if err != nil {
		devolverError(writer, "Unirse a Partida", err)
		return
	}

	// Comprobamos que la partida no esté completa
	if jugadores == maxJugadores {
		devolverError(writer, "Unirse a Partida", errors.New("No hay hueco en la partida"))
		return
	}

	// Comprobames que el usuario no esté participando en otra partida
	enPartida, err := dao.UsuarioEnPartida(globales.Db, &usuario)
	if err != nil {
		devolverError(writer, "Unirse a Partida", err)
		return
	}

	if enPartida {
		devolverError(writer, "Unirse a Partida", errors.New("El usuario ya está en otra partida"))
		return
	}

	publica, passwordHash, err := dao.ConsultarAcceso(globales.Db, &partida)
	if err != nil {
		devolverError(writer, "Unirse a Partida", err)
		return
	}

	if !publica {
		// Comprobamos que la contraseña sea correcta
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		if err != nil {
			devolverError(writer, "Unirse a Partida", errors.New("La contraseña no es correcta"))
			return
		}
	}

	// Else -> no está completa, el usuario no está en otra partida y la partida es pública o la contraseña es correcta
	err = dao.UnirseAPartida(globales.Db, &usuario, &partida)
	if err != nil {
		devolverError(writer, "Unirse a Partida", err)
		return
	}

	// Además, se aumenta el número de jugadores presente en ella, en su estado
	/*err = dao.ObtenerEstadoSerializado(globales.Db, &partida)
	if err != nil {
		devolverError(writer, "Unirse a Partida", err)
		return
	}

	partida.Estado.NumeroJugadores = partida.Estado.NumeroJugadores + 1*/

	devolverExito(writer)
}

// ObtenerPartidas devuelve un listado de partidas codificado en JSON, con el siguiente orden:
//	1- partidas privadas, de más a menos amigos presentes
//	2- partidas públicas, de más a menos amigos presentes
//	3- partidas públicas sin amigos: de más a menos jugadores
//	4- partidas privadas sin amigos: de más a menos jugadores
func ObtenerPartidas(writer http.ResponseWriter, request *http.Request) {
	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	amigos, err := dao.ObtenerAmigos(globales.Db, &usuario)
	if err != nil {
		devolverError(writer, "ObtenerPartidas", err)
	}

	log.Println("amigos del usuario:", amigos)

	partidas, err := dao.ObtenerPartidas(globales.Db)
	if err != nil {
		devolverError(writer, "ObtenerPartidas", err)
	}

	for _, p := range partidas {
		log.Println("Amigos en partida", p.IdPartida, ":", vo.ContarAmigos(amigos, p))
	}

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

	// Ordena partidas privadas de más a menos amigos
	sort.SliceStable(partidasPrivadas, func(i, j int) bool {
		// Orden: > a <
		return vo.ContarAmigos(amigos, partidasPrivadas[i]) > vo.ContarAmigos(amigos, partidasPrivadas[j])
	})

	// Ordena partidas públicas de más a menos amigos
	sort.SliceStable(partidasPublicas, func(i, j int) bool {
		// Orden: > a <
		return vo.ContarAmigos(amigos, partidasPublicas[i]) > vo.ContarAmigos(amigos, partidasPublicas[j])
	})

	// Extrae las partidas privadas sin amigos del usuario del slice y deja las partidas privadas con amigos
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

	// Ordena partidas privadas sin amigos de más a menos jugadores
	sort.SliceStable(partidasPrivadasSinAmigos, func(i, j int) bool {
		// Orden: > a <
		numeroJugadoresI, _, err1 := dao.ConsultarNumeroJugadores(globales.Db, &partidasPrivadasSinAmigos[i])
		numeroJugadoresJ, _, err2 := dao.ConsultarNumeroJugadores(globales.Db, &partidasPrivadasSinAmigos[i])

		if err1 != nil || err2 != nil {
			devolverError(writer, "ObtenerPartidas", err)
		}

		return numeroJugadoresI > numeroJugadoresJ
	})

	// Extrae las partidas públicas sin amigos del usuario del slice y deja las partidas públicas con amigos
	var partidasPublicasConAmigos []vo.Partida
	var partidasPublicasSinAmigos []vo.Partida
	for _, partida := range partidasPublicas {
		// Se ha llegado al punto en el slice a partir del cual no hay amigos
		if vo.ContarAmigos(amigos, partida) == 0 {
			partidasPublicasSinAmigos = append(partidasPublicasSinAmigos, partida)
		} else {
			partidasPublicasConAmigos = append(partidasPublicasConAmigos, partida)
		}
	}

	// Ordena partidas públicas sin amigos de más a menos jugadores
	sort.SliceStable(partidasPublicasSinAmigos, func(i, j int) bool {
		// Orden: > a <
		numeroJugadoresI, _, err1 := dao.ConsultarNumeroJugadores(globales.Db, &partidasPublicasSinAmigos[i])
		numeroJugadoresJ, _, err2 := dao.ConsultarNumeroJugadores(globales.Db, &partidasPublicasSinAmigos[i])

		if err1 != nil || err2 != nil {
			devolverError(writer, "ObtenerPartidas", err)
		}

		return numeroJugadoresI > numeroJugadoresJ
	})

	// Junta todos los slices, en orden
	var partidasOrdenadas []vo.Partida
	partidasOrdenadas = append(partidasOrdenadas, partidasPrivadasConAmigos...)
	partidasOrdenadas = append(partidasOrdenadas, partidasPublicasConAmigos...)
	partidasOrdenadas = append(partidasOrdenadas, partidasPublicasSinAmigos...)
	partidasOrdenadas = append(partidasOrdenadas, partidasPrivadasSinAmigos...)

	if err != nil {
		devolverError(writer, "ObtenerPartidas", err)
	} else {
		var elementos []vo.ElementoListaPartidas

		for _, p := range partidasOrdenadas {
			elementos = append(elementos, transformarAElementoListaPartidas(&p, amigos))
		}

		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(elementos)
	}
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
