package dao

import (
	"backend/vo"
	"bytes"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// CrearPartida crea una nueva partida, la cual será añadida a la base de datos
// El usuario especificará el número de jugadores máximos, si la partida es pública
// o privada, y la contraseña cuando sea necesario
func CrearPartida(db *sql.DB, usuario *vo.Usuario, partida *vo.Partida) (err error) {
	var IdPartida int
	password := sql.NullString{String: partida.PasswordHash, Valid: len(partida.PasswordHash) > 0}

	var estado bytes.Buffer
	encoder := gob.NewEncoder(&estado)
	err = encoder.Encode(partida.Estado)

	log.Println("partida.Estado antes:", partida.Estado)
	var estadoPartida []byte
	buf := bytes.NewBuffer(estadoPartida)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&partida.Estado)
	log.Println("partida.Estado despues:", partida.Estado)
	log.Println("bytes:", hex.EncodeToString(estado.Bytes()))

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
	return err
}

// UnirseAPartida crea una nueva entrada en la tabla "Participa" indicando que el usuario
// forma parte de la partida
func UnirseAPartida(db *sql.DB, usuario *vo.Usuario, partida *vo.Partida) (err error) {
	_, err = db.Exec(`INSERT INTO backend."Participa"("ID_partida", "nombreUsuario") VALUES($1, $2)`,
		partida.IdPartida, usuario.NombreUsuario)
	return err
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

// ConsultarNumeroJugadores devuelve el número actual de jugadores de una partida, además
// del número máximo de jugadores permitidos
func ConsultarNumeroJugadores(db *sql.DB, partida *vo.Partida) (jugadores, maxJugadores int, err error) {
	err = db.QueryRow(`SELECT "backend"."Partida"."maxJugadores" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&maxJugadores)
	err = db.QueryRow(`SELECT count(*) from (select distinct * from backend."Participa" 
		WHERE "ID_partida" = $1) AS sq`, partida.IdPartida).Scan(&jugadores)

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

// ObtenerPartidas devuelve un listado de todas las partidas que no están en curso almacenadas,
// ordenadas de privadas a públicas.
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

			partida.Jugadores = append(partida.Jugadores, vo.Usuario{NombreUsuario: nombre})
		}

		// Una vez escaneadas las columnas en los campos del struct, se obtiene el resto de campos no directos
		partida.PasswordHash = passwordHash.String

		log.Println("bytes en obtener:", hex.EncodeToString(estadoPartida))

		var estadoMio vo.EstadoPartida
		buf := bytes.NewBuffer(estadoPartida)
		log.Println("buffer:", hex.EncodeToString(buf.Bytes()))
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(&estadoMio)
		if err != nil {
			log.Println("que ha pasado", err)
			return partidas, err
		} else {
			log.Println("funciona")
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
