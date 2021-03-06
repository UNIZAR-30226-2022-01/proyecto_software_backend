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
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

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
// Tipo: POST
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
// [{carta}, {carta}, ...]
//
// Por ejemplo:
// [
//      	{
//        		"IdCarta": 1,
//        		"Tipo": 0,
//        		"Region": 29,
//        		"EsComodin": false
//        	},
//        	{
//        		"IdCarta": 20,
//        		"Tipo": 1,
//        		"Region": 22,
//        		"EsComodin": false
//        	}
// ]
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

// PasarDeFase permite al jugador actual cambiar de fase dentro de su propio turno, siendo estas fases Refuerzo,
// ataque y fortificación. Cada fase tendrá unas condiciones especiales para el cambio de turno:
// En el refuerzo, no podrá cambiar de fase si tiene más de 4 cartas o si le quedan tropas por asignar
// En el ataque, no podrá cambiar de fase si tiene más de 4 cartas o si tiene que ocupar un territorio y aún no lo ha hecho.
// En la fortificación podrá cambiar de fase (dándole el turno a otro jugador) libremente
//
// Si no es el turno del jugador, no está en una partida o no se cumplen las condiciones para el cambio de fase, devolverá
// un status 500 junto a un mensaje de error en el cuerpo, en otro caso devolverá status 200.
//
// Ruta: /api/pasarDeFase
// Tipo: POST
func PasarDeFase(writer http.ResponseWriter, request *http.Request) {
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
	err = partida.Estado.FinDeFase(usuario.NombreUsuario)
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

// Fortificar permite al jugador mover un número determinado de tropas de un territorio dado a otro.
//
// Las condiciones para poder fortificar son las siguientes:
// 		Ambos territorios pertenecen al jugador y son diferentes
//		Ambos territorios se encuentran conectados por algún camino en el mapa que cruce únicamente
//		territorios controlados por dicho jugador
//		El número de tropas del territorio origen debe ser mayor que 1
//		El número de tropas a mover es un número comprendido entre [1, num_tropas_territorio_1 - 1], de tal
//		forma que no se puede dejar el territorio origen sin tropas
//
// Si no es el turno del jugador, no está en una partida o no se cumplen las condiciones para la fortificación,
// se devolverá un status 500 junto a un mensaje de error en el cuerpo, en otro caso devolverá status 200
// y generará una acción de fortificación.
//
// Ruta: /api/fortificar/{id_territorio_origen}/{id_territorio_destino}/{num_tropas}
// Tipo: POST
func Fortificar(writer http.ResponseWriter, request *http.Request) {
	idTerritorioOrigen, err1 := strconv.Atoi(chi.URLParam(request, "id_territorio_origen"))
	idTerritorioDestino, err2 := strconv.Atoi(chi.URLParam(request, "id_territorio_destino"))
	numTropas, err3 := strconv.Atoi(chi.URLParam(request, "num_tropas"))

	if err1 != nil || err2 != nil || err3 != nil || numTropas == 0 {
		devolverError(writer, errors.New("Los identificadores de región y el número de tropas deben ser números naturales"))
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

	err = partida.Estado.FortificarTerritorio(idTerritorioOrigen, idTerritorioDestino, numTropas, usuario.NombreUsuario)
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

// Atacar permite a un usuario atacar a un territorio adyacente a alguna de sus regiones. Para ello, debe de tener por
// lo menos dos ejércitos en el territorio desde el que ataca. Al atacar deberá elegir el número de dados a lanzar, entre
// 1 y 3. Cabe destacar que será necesario tener al menos un ejército más que el número de dados a lanzar, por ejemplo, si
// quiero lanzar 3 dados, el territorio tendrá que tener 4 ejércitos por lo menos.
// Por otro lado el defensor tirará 2 dados si tiene 2 ejércitos o más, o 1 en el caso contrario.
//
// Para calcular el resultado del ataque, se compararán los dados con mayor valor de ambos jugadores. Si el atacante consigue
// un resultado mayor, el territorio defensor perderá una tropa. Por otro lado, si empatan o gana el defensor el territorio
// atacante perderá un ejército. En caso de que ambos jugadores hayan lanzado más de un dado, se repetirá el mismo proceso
// comparando el valor del segundo dado más alto de cada uno
//
// No se puede atacar en los siguientes casos: no es el turno deñ jugador, no es la fase de ataque, el jugador tiene más de
// 4 cartas, hay algún territorio sin ocupar, el territorio atacado no es adyacente, el territorio atacado no es de un rival,
// el número de dados no está entre 1 y 3 o el número de ejércitos no supera el número de dados.
//
// Ruta: /api/atacar/{id_territorio_origen}/{id_territorio_destino}/{num_dados}
// Tipo: POST
func Atacar(writer http.ResponseWriter, request *http.Request) {
	origen, err1 := strconv.Atoi(chi.URLParam(request, "id_territorio_origen"))
	destino, err2 := strconv.Atoi(chi.URLParam(request, "id_territorio_destino"))
	numDados, err3 := strconv.Atoi(chi.URLParam(request, "num_dados"))

	if err1 != nil || err2 != nil || err3 != nil {
		devolverError(writer, errors.New("Los identificadores de región y el número de tropas deben ser números naturales"))
		return
	}

	if origen < 0 || origen >= logica_juego.NUM_REGIONES || destino < 0 || destino >= logica_juego.NUM_REGIONES {
		devolverError(writer, errors.New("El identificador de región debe estar entre 0 y"+
			strconv.Itoa(logica_juego.NUM_REGIONES)+"no incluido"))
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
	err = partida.Estado.Ataque(logica_juego.NumRegion(origen), logica_juego.NumRegion(destino), numDados, usuario.NombreUsuario)
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

// Ocupar permite a un usuario ocupar un territorio sin tropas, especificando el territorio a ocupar y el número de
// tropas que quiere mover a él. Dichas tropas se moverán desde el territorio con el que conquistó la región a ser ocupada.
//
// Para ocupar se deben cumplir las siguientes condiciones: hay alguna región sin tropas, dicha región es adyacente a
// la región desde la que se inició el último ataque, la ocupación se realiza durante el turno del jugador y en la fase
// de ataque, el número de tropas asignadas por la ocupación no deja al territorio origen sin tropas, el número de tropas
// asignadas es mayor al número de dados usados en el último ataque menos el número de ejércitos que perdió el atacante
// en dicho ataque.
//
// Cabe destacar que siempre que un territorio quede sin tropas tras un ataque, el juego no permitirá continuar atacando
// ni cambiar de fase o turno hasta que dicho territorio sea ocupado, de manera que solo podrá haber un territorio
// sin ocupar a la vez.
// Ruta: /api/ocupar/{territorio_a_ocupar}/{num_ejercitos}
// Tipo: POST
func Ocupar(writer http.ResponseWriter, request *http.Request) {
	regionAOcupar, err1 := strconv.Atoi(chi.URLParam(request, "territorio_a_ocupar"))
	numEjercitos, err2 := strconv.Atoi(chi.URLParam(request, "num_ejercitos"))

	if err1 != nil || err2 != nil {
		devolverError(writer, errors.New("Los identificadores de región y el número de tropas deben ser números naturales"))
		return
	}

	if regionAOcupar < 0 || regionAOcupar >= logica_juego.NUM_REGIONES {
		devolverError(writer, errors.New("El identificador de región debe estar entre 0 y"+
			strconv.Itoa(logica_juego.NUM_REGIONES)+"no incluido"))
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
	err = partida.Estado.Ocupar(logica_juego.NumRegion(regionAOcupar), numEjercitos, usuario.NombreUsuario)
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
