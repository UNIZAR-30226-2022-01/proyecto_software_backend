package dao

import (
	"backend/vo"
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"
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
		log.Println("error en select de consultar cookie:", err)
	}
	b.Write(bytearray)
	err = decoder.Decode(&cookie)

	if err != nil {
		log.Println("error al serializar cookie:", err)
	}

	return cookie, err
}

// InsertarCookie registra una cookie para el usuario dado, buscando por su nombre. En caso de fallo o no
// encontrarse, devuelve un error.
func InsertarCookie(db *sql.DB, usuario *vo.Usuario) (err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(usuario.CookieSesion)

	if err != nil {
		log.Println("error al serializar cookie")
	}

	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "cookieSesion" = $1 WHERE "backend"."Usuario"."nombreUsuario" = $2`, b.Bytes(), usuario.NombreUsuario)

	if err != nil {
		log.Println("error en update de insertar cookie:", err)
	}

	return err
}

// InsertarUsuario registra un usuario dados sus datos. En caso de fallo, devuelve un error.
func InsertarUsuario(db *sql.DB, usuario *vo.Usuario) (err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(usuario.CookieSesion)

	_, err = db.Exec(`INSERT INTO "backend"."Usuario"("email", "nombreUsuario", "passwordHash", 
		"biografia", "cookieSesion", "puntos", "partidasGanadas", "partidasTotales", "ID_dado", "ID_ficha")
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, usuario.Email, usuario.NombreUsuario, usuario.PasswordHash,
		usuario.Biografia, b.Bytes(), usuario.Puntos, usuario.PartidasGanadas, usuario.PartidasTotales, usuario.ID_dado, usuario.ID_ficha)

	return err
}

// ConsultarPasswordHash devuelve el hash de contraseña del usuario dado, buscando por su nombre. En caso de fallo o no
// encontrarse, devuelve un error.
func ConsultarPasswordHash(db *sql.DB, usuario *vo.Usuario) (hash string, err error) {
	err = db.QueryRow(`SELECT "backend"."Usuario"."passwordHash" FROM "backend"."Usuario"
		WHERE "backend"."Usuario"."nombreUsuario" = $1`, usuario.NombreUsuario).Scan(&hash)

	if err != nil {
		log.Println("error en select de consultar password:", err)
	}

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
	log.Println("Aceptando de", emisor.NombreUsuario, "a", receptor.NombreUsuario)

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

// ObtenerAmigos devuelve una lista de usuarios (con su nombre de usuario rellenado)
// que son amigos del usuario indicado, o error en caso de fallo.
func ObtenerAmigos(db *sql.DB, usuario *vo.Usuario) (amigos []vo.Usuario, err error) {
	rows, err := db.Query(`SELECT backend."EsAmigo"."nombreUsuario1", backend."EsAmigo"."nombreUsuario2"
									FROM backend."EsAmigo" 
									WHERE $1 in (backend."EsAmigo"."nombreUsuario1", backend."EsAmigo"."nombreUsuario2" )`, usuario.NombreUsuario)
	defer rows.Close()
	if err != nil {
		log.Println("Error al consultar filas en ObtenerPartidas:", err)
		return amigos, err
	}

	for rows.Next() {
		var amigo vo.Usuario
		var nombre1 string
		var nombre2 string
		err = rows.Scan(&nombre1, &nombre2)

		if err != nil {
			log.Println("Error al recuperar fila en ObtenerPartidas:", err)
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
