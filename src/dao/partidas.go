package dao

import (
	"backend/vo"
	"database/sql"
	"log"
)

// CrearPartida crea una nueva partida, la cual será añadida a la base de datos
// EL usuario especificará el número de jugadores máximos, si la partida es pública
// o privada, y la contrasña cuando sea necesario
func CrearPartida(db *sql.DB, usuario *vo.Usuario, partida *vo.Partida) (err error) {
	var IdPartida int
	password := sql.NullString{String: partida.PasswordHash, Valid: len(partida.PasswordHash) > 0}
	err = db.QueryRow(`INSERT INTO "backend"."Partida"("estadoPartida", "mensajes", "esPublica",
        "passwordHash", "enCurso", "maxJugadores") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
		partida.Estado, partida.Mensajes, partida.EsPublica, password, partida.EnCurso, partida.MaxNumeroJugadores).Scan(&IdPartida)
	if err != nil {
		return err
	}
	log.Println("El id de la partida introducida es: ", IdPartida)

	_, err = db.Exec(`INSERT INTO "backend"."Participa"("ID_partida", "nombreUsuario") 
					VALUES ($1, $2)`, IdPartida, usuario.NombreUsuario)
	return err
}

// UnisrseAPartida crea una nueva entrada en la tabla "Participa" indicando que el usuario
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
	err = db.QueryRow(`SELECT "backend"."Partida"."passwordHash" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&hash)
	err = db.QueryRow(`SELECT "backend"."Partida"."esPublica" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&esPublica)
	if err != nil {
		log.Println("error en select de consultar acceso a partida:", err)
	}

	return esPublica, hash, err
}

// ConsultarNumeroJugadores devuelve el número actual de jugadores de una partida, además
// del número máximo de jugadores permitidos
func ConsultarNumeroJugadores(db *sql.DB, partida *vo.Partida) (jugadores, maxJugadores int, err error) {
	err = db.QueryRow(`SELECT "backend"."Partida"."maxJugadores" FROM "backend"."Partida"
		WHERE "backend"."Partida"."id" = $1`, partida.IdPartida).Scan(&maxJugadores)
	err = db.QueryRow(`SELECT count(*) from (select distinct * from backend."Participa" 
		WHERE "ID_partida" = $1) AS sq`, partida.IdPartida).Scan(&jugadores)
	if err != nil {
		log.Println("error en select de consultar número de jugadores:", err)
	}
	log.Println("Número de jugadores:", jugadores, ",jugadores maximos: ", maxJugadores, "en la partida", partida.IdPartida)
	return jugadores, maxJugadores, err
}
