package handlers

import (
	"backend/dao"
	"backend/globales"
	"backend/middleware"
	"backend/vo"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
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
	partida = vo.Partida{0, esPublica, partida.PasswordHash, false, 1, maxJugadores, aux, aux}

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

	devolverExito(writer)
}
