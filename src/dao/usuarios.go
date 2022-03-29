package dao

import (
	"backend/vo"
	"bytes"
	"database/sql"
	"encoding/gob"
	"net/http"
)

// ConsultarCookie devuelve una cookie del usuario dado, buscando por su nombre. En caso de fallo o no
// encontrarse, devuelve un error.
func ConsultarCookie(db *sql.DB, usuario *vo.Usuario) (cookie http.Cookie, err error) {
	var b bytes.Buffer
	decoder := gob.NewDecoder(&b)

	var bytearray []byte

	err = db.QueryRow(`SELECT "backend"."Usuario"."cookieSesion" FROM "backend"."Usuario"
		WHERE "backend"."Usuario"."nombreUsuario" = $1`, usuario.NombreUsuario).Scan(&bytearray)
	if err != nil {
		return cookie, err
	}

	b.Write(bytearray)
	err = decoder.Decode(&cookie)

	return cookie, err
}

// InsertarCookie registra una cookie para el usuario dado, buscando por su nombre. En caso de fallo o no
// encontrarse, devuelve un error.
func InsertarCookie(db *sql.DB, usuario *vo.Usuario) (err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(usuario.CookieSesion)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "cookieSesion" = $1 WHERE "backend"."Usuario"."nombreUsuario" = $2`, b.Bytes(), usuario.NombreUsuario)

	return err
}

// InsertarUsuario registra un usuario dados sus datos. En caso de fallo, devuelve un error.
func InsertarUsuario(db *sql.DB, usuario *vo.Usuario) (err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(usuario.CookieSesion)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO "backend"."Usuario"("email", "nombreUsuario", "passwordHash", 
		"biografia", "cookieSesion", "puntos", "partidasGanadas", "partidasTotales", "ID_dado", "ID_ficha")
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, usuario.Email, usuario.NombreUsuario, usuario.PasswordHash,
		usuario.Biografia, b.Bytes(), usuario.Puntos, usuario.PartidasGanadas, usuario.PartidasTotales, usuario.ID_dado, usuario.ID_ficha)

	return err
}

func ObtenerUsuario(db *sql.DB, nombreUsuario string) (usuario vo.Usuario, err error) {
	var b bytes.Buffer
	decoder := gob.NewDecoder(&b)

	var bytearray []byte
	err = db.QueryRow(`SELECT "email", "nombreUsuario", "passwordHash", "biografia", "cookieSesion", 
		"partidasGanadas", "partidasTotales", "puntos", "ID_dado", "ID_ficha" 
		FROM backend."Usuario" WHERE backend."Usuario"."nombreUsuario" = $1`, nombreUsuario).Scan(
		&usuario.Email, &usuario.NombreUsuario, &usuario.PasswordHash, &usuario.Biografia, &bytearray,
		&usuario.PartidasGanadas, &usuario.PartidasTotales, &usuario.Puntos, &usuario.ID_dado, &usuario.ID_ficha)
	b.Write(bytearray)
	err = decoder.Decode(&usuario.CookieSesion)
	return usuario, err
}

// ConsultarPasswordHash devuelve el hash de contraseña del usuario dado, buscando por su nombre. En caso de fallo o no
// encontrarse, devuelve un error.
func ConsultarPasswordHash(db *sql.DB, usuario *vo.Usuario) (hash string, err error) {
	err = db.QueryRow(`SELECT "backend"."Usuario"."passwordHash" FROM "backend"."Usuario"
		WHERE "backend"."Usuario"."nombreUsuario" = $1`, usuario.NombreUsuario).Scan(&hash)
	return hash, err
}

// CrearSolicitudAmistad registra una solicitud de amistad entre los usuarios emisor y receptor. En caso de fallo o no
// encontrarse alguno de ellos, devuelve un error.
func CrearSolicitudAmistad(db *sql.DB, emisor *vo.Usuario, receptor *vo.Usuario) error {
	_, err := db.Exec(`INSERT INTO "backend"."EsAmigo"("nombreUsuario1", "nombreUsuario2", "pendiente") VALUES($1, $2, $3)`,
		emisor.NombreUsuario, receptor.NombreUsuario, true)

	return err
}

// AceptarSolicitudAmistad registra una solicitud de amistad existente como aceptada entre los usuarios emisor y receptor.
// En caso de fallo o no encontrarse alguno de ellos o la solicitud, devuelve un error.
func AceptarSolicitudAmistad(db *sql.DB, emisor *vo.Usuario, receptor *vo.Usuario) error {
	_, err := db.Exec(`UPDATE "backend"."EsAmigo" SET "pendiente" = false WHERE "nombreUsuario1" = $1 AND "nombreUsuario2" = $2`,
		receptor.NombreUsuario, emisor.NombreUsuario)
	return err
}

// AceptarSolicitudAmistad elimina una solicitud de amistad existente entre los usuarios emisor y receptor.
// En caso de fallo o no encontrarse alguno de ellos o la solicitud, devuelve un error.
func RechazarSolicitudAmistad(db *sql.DB, emisor *vo.Usuario, receptor *vo.Usuario) error {
	_, err := db.Exec(`DELETE FROM "backend"."EsAmigo" WHERE "nombreUsuario1" = $1 AND "nombreUsuario2" = $2`,
		receptor.NombreUsuario, emisor.NombreUsuario)
	return err
}

// UsuarioEnPartida devolverá true en caso de que un usuario ya esté participando en una partida
func UsuarioEnPartida(db *sql.DB, usuario *vo.Usuario) (EnPartida bool, err error) {
	err = db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."Participa" WHERE "nombreUsuario" = $1)`, usuario.NombreUsuario).Scan(&EnPartida)
	return EnPartida, err
}

// PartidaUsuario devuelve el ID de la partida en la que participa un usuario, y error en cualquier otro caso
func PartidaUsuario(db *sql.DB, usuario *vo.Usuario) (idPartida int, err error) {
	err = db.QueryRow(`SELECT backend."Participa"."ID_partida"  FROM backend."Participa" WHERE "nombreUsuario" = $1`, usuario.NombreUsuario).Scan(&idPartida)

	return idPartida, err
}

// ObtenerAmigos devuelve una lista de usuarios (con su nombre de usuario rellenado)
// que son amigos del usuario indicado, o error en caso de fallo.
func ObtenerAmigos(db *sql.DB, usuario *vo.Usuario) (amigos []vo.Usuario, err error) {
	rows, err := db.Query(`SELECT backend."EsAmigo"."nombreUsuario1", backend."EsAmigo"."nombreUsuario2"
									FROM backend."EsAmigo" 
									WHERE $1 in (backend."EsAmigo"."nombreUsuario1", backend."EsAmigo"."nombreUsuario2" )`, usuario.NombreUsuario)
	defer rows.Close()
	if err != nil {
		return amigos, err
	}

	for rows.Next() {
		var amigo vo.Usuario
		var nombre1 string
		var nombre2 string
		err = rows.Scan(&nombre1, &nombre2)

		if err != nil {
			return amigos, err
		}

		// Elige el nombre de la tupla que no coincide con el del usuario
		if nombre1 != usuario.NombreUsuario {
			amigo = vo.Usuario{NombreUsuario: nombre1}
		} else {
			amigo = vo.Usuario{NombreUsuario: nombre2}
		}
		amigos = append(amigos, amigo)
	}

	return amigos, nil
}

// ObtenerUsuariosSimilares devuelve el nombre de usuario de todos los usuarios registrados cuyo nombre sea similar
// a uno indicado, ordenados alfabéticamente
func ObtenerUsuariosSimilares(db *sql.DB, nombre string) (usuarios []string, err error) {
	patron := nombre + "%"
	rows, err := db.Query(`SELECT backend."Usuario"."nombreUsuario" FROM backend."Usuario" 
		WHERE "nombreUsuario" LIKE $1 ORDER BY backend."Usuario"."nombreUsuario" ASC `, patron)
	if err != nil {
		return usuarios, err
	}
	defer rows.Close()
	for rows.Next() {
		var usuario string
		err = rows.Scan(&usuario)
		if err != nil {
			return usuarios, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, err
}

// ExisteUsuario devuelve true si hay algún usuario con el nombre "nombre" registrado
func ExisteUsuario(db *sql.DB, nombre string) bool {
	var existe bool
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."Usuario"
		WHERE "nombreUsuario" = $1)`, nombre).Scan(existe)
	if err != nil {
		return false
	}
	return existe
}

// ExisteEmail devuelve true si hay algún usuario con el email "email" registrado
func ExisteEmail(db *sql.DB, email string) bool {
	var existe bool
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."Usuario"
		WHERE "email" = $1)`, email).Scan(existe)
	if err != nil {
		return false
	}
	return existe
}
