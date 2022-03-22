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

// CrearPartida crea una nueva partida, para la que se definirá el número máxmio de jugadores,
// si es pública o privada, y la contraseña en caso de que fuera necesario
func CrearPartida(writer http.ResponseWriter, request *http.Request) {
	password := request.FormValue("password")
	maxJugadores, err := strconv.Atoi(request.FormValue("maxJugadores"))
	nombreUsuario := middleware.ObtenerUsuarioCookie(request)
	tipoPartida := request.FormValue("tipo")
	esPublica := tipoPartida == "Publica"

	if err != nil {
		devolverError(writer, "Crear Partida", err)
		return
	}
	if maxJugadores < 2 || maxJugadores > 6 {
		devolverError(writer, "Crear Partida", errors.New("El número de jugadores debe estar entre 2 y 6"))
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

	aux := 1 // TODO, para que mensaje y estado no den error por nulo
	usuario := vo.Usuario{"", nombreUsuario, "", "", http.Cookie{}, 0, 0, 0, 0, 0}
	partida = vo.Partida{0, esPublica, partida.PasswordHash, false, 1, maxJugadores, nil, aux, aux}

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
//
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

	devolverExito(writer)
}

// ObtenerPartidas devuelve un listado de partidas codificado en JSON, con el siguiente orden:
//	1- partidas privadas, de más a menos amigos presentes
//	2- partidas públicas, de más a menos amigos presentes
//	3- partidas públicas sin amigos: de más a menos jugadores
func ObtenerPartidas(writer http.ResponseWriter, request *http.Request) {
	usuario := vo.Usuario{NombreUsuario: middleware.ObtenerUsuarioCookie(request)}

	amigos, err := dao.ObtenerAmigos(globales.Db, &usuario)
	if err != nil {
		// TODO en json
		devolverError(writer, "ObtenerPartidas", err)
	}

	partidas, err := dao.ObtenerPartidas(globales.Db)
	if err != nil {
		// TODO en json
		devolverError(writer, "ObtenerPartidas", err)
	}

	// Extrae las partidas privadas del slice y deja las partidas públicas
	var partidasPrivadas []vo.Partida
	for i, partida := range partidas {
		if !partida.EsPublica {
			partidasPrivadas = append(partidasPrivadas, partida)
			partidas = append(partidas[:i], partidas[i+1:]...)
		} else {
			break
		}
	}

	// Ordena partidas privadas de más a menos amigos
	sort.SliceStable(partidasPrivadas, func(i, j int) bool {
		// Orden: > a <
		return vo.ContarAmigos(amigos, partidasPrivadas[i]) > vo.ContarAmigos(amigos, partidasPrivadas[j])
	})

	// Ordena partidas públicas de más a menos amigos
	sort.SliceStable(partidas, func(i, j int) bool {
		// Orden: > a <
		return vo.ContarAmigos(amigos, partidas[i]) > vo.ContarAmigos(amigos, partidas[j])
	})

	// Extrae las partidas públicas sin amigos del usuario del slice y deja las partidas públicas con amigos
	var partidasPublicasSinAmigos []vo.Partida
	for i, partida := range partidas {
		// Se ha llegado al punto en el slice a partir del cual no hay amigos
		if vo.ContarAmigos(amigos, partida) == 0 {
			partidasPublicasSinAmigos = partidas[i:]
			partidas = partidas[:i]
			break
		}
	}

	// Ordena partidas públicas sin amigos de más a menos jugadores
	sort.SliceStable(partidasPublicasSinAmigos, func(i, j int) bool {
		// Orden: > a <
		return partidasPublicasSinAmigos[i].NumeroJugadores > partidasPublicasSinAmigos[j].NumeroJugadores
	})

	/*log.Println("partidas privadas despues de sort:", partidasPrivadas)
	log.Println("partidas públicas despues de sort:", partidas)
	log.Println("partidas públicas sin amigos despues de sort:", partidasPublicasSinAmigos)*/

	// Junta todos los slices, en orden
	partidasPrivadas = append(partidasPrivadas, partidas...)
	partidasPrivadas = append(partidasPrivadas, partidasPublicasSinAmigos...)

	log.Println("partidas ordenadas al final:", partidasPrivadas)

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(partidasPrivadas)
	if err != nil {
		// TODO en json
		devolverError(writer, "ObtenerPartidas", err)
	}
}
