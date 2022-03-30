package dao

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
)

// CrearPartida crea una nueva partida, la cual será añadida a la base de datos
// El usuario especificará el número de jugadores máximos, si la partida es pública
// o privada, y la contraseña cuando sea necesario
// Modifica el objeto partida utilizando, especificando su identificador asignado
// al almacenarla en la base de datos
func CrearPartida(db *sql.DB, usuario *vo.Usuario, partida *vo.Partida) (err error) {
	var IdPartida int
	password := sql.NullString{String: partida.PasswordHash, Valid: len(partida.PasswordHash) > 0}

	var estado bytes.Buffer
	encoder := gob.NewEncoder(&estado)
	err = encoder.Encode(partida.Estado)

	var mensajes bytes.Buffer
	encoder = gob.NewEncoder(&mensajes)
	err = encoder.Encode(partida.Mensajes)

	err = db.QueryRow(`INSERT INTO "backend"."Partida"("estadoPartida", "mensajes", "esPublica",
        "passwordHash", "enCurso", "maxJugadores") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
		estado.Bytes(), mensajes.Bytes(), partida.EsPublica, password, partida.EnCurso, partida.MaxNumeroJugadores).Scan(&IdPartida)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO "backend"."Participa"("ID_partida", "nombreUsuario") 
					VALUES ($1, $2)`, IdPartida, usuario.NombreUsuario)
	partida.IdPartida = IdPartida

	return err
}

// BorrarPartida borra la partida indicada si existe, o devuelve un error en caso contrario.
func BorrarPartida(db *sql.DB, partida *vo.Partida) error {
	_, err := db.Exec(`DELETE FROM backend."Partida" WHERE backend."Partida".id= $1;`, partida.IdPartida)
	return err
}

// UnirseAPartida crea una nueva entrada en la tabla "Participa" indicando que el usuario
// forma parte de la partida
func UnirseAPartida(db *sql.DB, usuario *vo.Usuario, partida *vo.Partida) (err error) {
	_, err = db.Exec(`INSERT INTO backend."Participa"("ID_partida", "nombreUsuario") VALUES($1, $2)`,
		partida.IdPartida, usuario.NombreUsuario)
	return err
}

// AbandonarLobby intenta abandonar una partida dada si no está en curso, o devuelve un error apropiado en caso contrario ya formateado.
// Adicionalmente, si la partida se queda sin jugadores, se borrará.
func AbandonarLobby(db *sql.DB, usuario *vo.Usuario) (err error) {
	idPartida := 0

	err = db.QueryRow(`SELECT "backend"."Participa"."ID_partida" FROM "backend"."Participa" WHERE "backend"."Participa"."nombreUsuario" = $1`, usuario.NombreUsuario).Scan(&idPartida)
	if err != nil && err != sql.ErrNoRows { // Error de SQL general
		return errors.New("Se ha producido un error al procesar los datos.")
	} else if err == sql.ErrNoRows {
		return errors.New("No estás participando en ninguna partida.")
	}

	enCurso := false
	err = db.QueryRow(`SELECT "backend"."Partida"."enCurso" FROM "backend"."Partida" WHERE "backend"."Partida"."id" = $1`, idPartida).Scan(&enCurso)
	if err != nil { // Error de SQL general
		return errors.New("Se ha producido un error al procesar los datos.")
	} else if enCurso {
		return errors.New("La partida ya está en curso.")
	}

	// Si no, la partida no está en curso y está participando en ella
	_, err = db.Exec(`DELETE FROM backend."Participa" WHERE "backend"."Participa"."nombreUsuario" = $1`, usuario.NombreUsuario)
	if err != nil {
		return errors.New("Se ha producido un error al procesar los datos.")
	}

	// Se comprueba si la partida se ha quedado sin usuarios y, si lo está, se borra
	numUsuarios := 0
	err = db.QueryRow(`SELECT COUNT(*) FROM backend."Participa" where backend."Participa"."ID_partida" = $1;`, idPartida).Scan(&numUsuarios)
	if err != nil {
		return errors.New("Se ha producido un error al procesar los datos.")
	}

	if numUsuarios == 0 {
		err = BorrarPartida(db, &vo.Partida{IdPartida: idPartida})
		if err != nil {
			return errors.New("Se ha producido un error al procesar los datos.")
		} else {
			return nil
		}
	} else {
		return nil
	}
}

// ObtenerEstadoLobby devuelve el estado del lobby de una partida identificada por su id
// Devuelve si es pública o no, si está o no en curso, el número máximo de jugadores y
// los jugadores que se encuentran en el lobby
func ObtenerEstadoLobby(db *sql.DB, idPartida int) (estado vo.EstadoLobby, err error) {
	err = db.QueryRow(`SELECT "enCurso", "esPublica", "maxJugadores" 
		FROM backend."Partida" WHERE "Partida".id = $1`, idPartida).Scan(&estado.EnCurso,
		&estado.EsPublico, &estado.MaxJugadores)
	if err != nil {
		return estado, err
	}

	rows, err := db.Query(`SELECT backend."Participa"."nombreUsuario" FROM backend."Participa" WHERE
			"ID_partida" = $1 ORDER BY "nombreUsuario" ASC`, idPartida)
	defer rows.Close()
	if err != nil {
		return estado, err
	}

	for rows.Next() {
		var jugador string
		err = rows.Scan(&jugador)
		if err != nil {
			return estado, err
		}

		estado.NombresJugadores = append(estado.NombresJugadores, jugador)
	}

	estado.Jugadores = len(estado.NombresJugadores)
	return estado, err
}

// ConsultarAcceso devuelve los permisos de acceso de una partida en concreto
// El parámetro de salida "esPublica" indicará si la partida es pública o no
// El parámetro hash corresponderá al hash de la contraseña para el acceso a la partida
func ConsultarAcceso(db *sql.DB, partida *vo.Partida) (esPublica bool, hash string, err error) {
	// TODO cambiar por una sola consulta
	err = db.QueryRow(`SELECT "backend"."Partida"."passwordHash" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&hash)
	err = db.QueryRow(`SELECT "backend"."Partida"."esPublica" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&esPublica)

	return esPublica, hash, err
}

// ConsultarJugadoresPartida devuelve los jugadores de una partida, además
// del número máximo de jugadores permitidos
func ConsultarJugadoresPartida(db *sql.DB, partida *vo.Partida) (jugadores []vo.Usuario, maxJugadores int, err error) {
	err = db.QueryRow(`SELECT "backend"."Partida"."maxJugadores" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&maxJugadores)
	if err != nil {
		return jugadores, maxJugadores, err
	}
	rows, err := db.Query(`select "nombreUsuario" from backend."Participa" WHERE "ID_partida" = $1`, partida.IdPartida)
	defer rows.Close()

	if err != nil {
		return jugadores, maxJugadores, err
	}

	for rows.Next() {
		var jugador vo.Usuario
		err = rows.Scan(&jugador.NombreUsuario)
		if err != nil {
			return jugadores, maxJugadores, err
		}
		jugadores = append(jugadores, jugador)
	}

	return jugadores, maxJugadores, err
}

// AlmacenarEstadoSerializado almacena el estado una partida dada, serializado a bytes. Devuelve un error en fallo.
func AlmacenarEstadoSerializado(db *sql.DB, partida *vo.Partida) (err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(partida.Estado)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE "backend"."Partida" SET "estadoPartida" = $1 WHERE "backend"."Partida".id = $2`, b.Bytes(), partida.IdPartida)

	return err
}

// AlmacenarMensajes almacena los mensajes una partida dada, serializados a bytes. Devuelve un error en fallo.
func AlmacenarMensajes(db *sql.DB, partida *vo.Partida) (err error) {
	var mensajes bytes.Buffer
	encoder := gob.NewEncoder(&mensajes)
	err = encoder.Encode(partida.Mensajes)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE "backend"."Partida" SET "mensajes" = $1 WHERE "backend"."Partida".id = $2`, mensajes.Bytes(), partida.IdPartida)
	return err
}

// ObtenerEstadoSerializado obtiene una partida existente con el ID indicado y deserializa su estado en ella, o devuelve un error en fallo.
func ObtenerEstadoSerializado(db *sql.DB, partida *vo.Partida) (err error) {
	var estadoPartida []byte
	err = db.QueryRow(`SELECT "backend"."Partida"."estadoPartida" FROM "backend"."Partida" WHERE "backend"."Partida".id = $1`, partida.IdPartida).Scan(&estadoPartida)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(estadoPartida)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&partida.Estado)

	return err
}

// ObtenerMensajes obtiene los mensajes partida existente con el ID indicado y los deserializa en ella, o devuelve un error en fallo.
func ObtenerMensajes(db *sql.DB, partida *vo.Partida) (err error) {
	var mensajesPartida []byte
	err = db.QueryRow(`SELECT "backend"."Partida"."estadoPartida" FROM "backend"."Partida" WHERE "backend"."Partida".id = $1`, partida.IdPartida).Scan(&mensajesPartida)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(mensajesPartida)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&partida.Mensajes)

	return err
}

// ObtenerPartidas devuelve un listado de todas las partidas, ordenadas de privadas a públicas.
func ObtenerPartidas(db *sql.DB) (partidas []vo.Partida, err error) {
	// Ordena por defecto de false a true
	rows, err := db.Query(`SELECT id, "estadoPartida", mensajes, "esPublica", "passwordHash", "enCurso", "maxJugadores" FROM backend."Partida" order by backend."Partida"."esPublica"`)
	defer rows.Close()
	if err != nil {
		return partidas, err
	}

	for rows.Next() {
		var estadoPartida []byte
		var mensajes []byte
		var partida vo.Partida
		var passwordHash sql.NullString
		err = rows.Scan(&partida.IdPartida, &estadoPartida, &mensajes, &partida.EsPublica, &passwordHash, &partida.EnCurso, &partida.MaxNumeroJugadores)
		if err != nil {
			return partidas, err
		}

		// Una vez escaneadas las columnas en los campos del struct, se obtiene el resto de campos no directos
		partida.PasswordHash = passwordHash.String

		buf := bytes.NewBuffer(estadoPartida)
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(&partida.Estado)
		if err != nil {
			return partidas, err
		}

		buf = bytes.NewBuffer(mensajes)
		decoder = gob.NewDecoder(buf)
		err = decoder.Decode(&partida.Mensajes)
		if err != nil {
			return partidas, err
		} else {
			partidas = append(partidas, partida)
		}
	}

	return partidas, nil
}

// ObtenerPartida devuelve una partida, dado su id, o error en cualquier otro caso.
func ObtenerPartida(db *sql.DB, idP int) (partida vo.Partida, err error) {
	var estadoPartida []byte
	var mensajes []byte
	var passwordHash sql.NullString

	err = db.QueryRow(`SELECT id, "estadoPartida", mensajes, "esPublica", "passwordHash", "enCurso", "maxJugadores" FROM backend."Partida" WHERE backend."Partida".id = $1`, idP).Scan(
		&partida.IdPartida, &estadoPartida, &mensajes, &partida.EsPublica, &passwordHash, &partida.EnCurso, &partida.MaxNumeroJugadores)
	if err != nil {
		return partida, err
	}

	// Una vez escaneadas las columnas en los campos del struct, se obtiene el resto de campos no directos
	partida.PasswordHash = passwordHash.String

	buf := bytes.NewBuffer(estadoPartida)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&partida.Estado)
	if err != nil {
		return partida, err
	}

	buf = bytes.NewBuffer(mensajes)
	decoder = gob.NewDecoder(buf)
	err = decoder.Decode(&partida.Mensajes)

	return partida, err
}

// ObtenerPartidasNoEnCurso devuelve un listado de todas las partidas que no están en curso almacenadas, ordenadas de privadas a públicas.
func ObtenerPartidasNoEnCurso(db *sql.DB) (partidas []vo.Partida, err error) {
	// Ordena por defecto de false a true
	rows, err := db.Query(`SELECT id, "estadoPartida", mensajes, "esPublica", "passwordHash", "enCurso", "maxJugadores" FROM backend."Partida" WHERE "enCurso" = false ORDER BY backend."Partida"."esPublica"`)
	defer rows.Close()
	if err != nil {
		return partidas, err
	}

	for rows.Next() {
		var estadoPartida []byte
		var mensajes []byte
		var partida vo.Partida
		var passwordHash sql.NullString
		err = rows.Scan(&partida.IdPartida, &estadoPartida, &mensajes, &partida.EsPublica, &passwordHash, &partida.EnCurso, &partida.MaxNumeroJugadores)
		if err != nil {
			return partidas, err
		}

		// Ahora se obtienen los usuarios participantes, y su número
		rowsInternas, err := db.Query(`SELECT "nombreUsuario" FROM backend."Participa" WHERE "ID_partida"= $1`, partida.IdPartida)
		defer rowsInternas.Close()
		if err != nil {
			return partidas, err
		}

		for rowsInternas.Next() {
			var nombre string
			err = rowsInternas.Scan(&nombre)
			if err != nil {
				return partidas, err
			}
		}

		// Una vez escaneadas las columnas en los campos del struct, se obtiene el resto de campos no directos
		partida.PasswordHash = passwordHash.String

		buf := bytes.NewBuffer(estadoPartida)
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(&partida.Estado)
		if err != nil {
			return partidas, err
		}

		buf = bytes.NewBuffer(mensajes)
		decoder = gob.NewDecoder(buf)
		err = decoder.Decode(&partida.Mensajes)
		if err != nil {
			return partidas, err
		} else {
			partidas = append(partidas, partida)
		}
	}

	return partidas, nil
}

func EmpezarPartida(db *sql.DB, idP int) error {
	_, err := db.Exec(`UPDATE backend."Partida" SET "enCurso"=true WHERE id=$1`, idP)
	return err
}
