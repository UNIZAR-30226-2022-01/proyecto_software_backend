package dao

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
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
	var existe bool
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."EsAmigo" WHERE "nombreUsuario1" = $1 AND "nombreUsuario2" = $2)`,
		receptor.NombreUsuario, emisor.NombreUsuario).Scan(&existe)

	if err != nil {
		return err
	}

	if existe {
		return errors.New("No puedes enviar una solicitud de amistad a un amigo," +
			" o a alguien que te ha enviado una solicitud pendiente")
	}

	_, err = db.Exec(`INSERT INTO "backend"."EsAmigo"("nombreUsuario1", "nombreUsuario2", "pendiente") VALUES($1, $2, $3)`,
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

// ConsultarSolicitudesPendientes devuelve una lista en la que se indican los nombres de usuario
// que han enviado una solicitud de amistad a "usuario", estando dicha solicitud pendiente.
func ConsultarSolicitudesPendientes(db *sql.DB, usuario *vo.Usuario) (usuarios []string, err error) {
	rows, err := db.Query(`SELECT "nombreUsuario1" FROM backend."EsAmigo" 
		WHERE "nombreUsuario2" = $1 AND pendiente ORDER BY "nombreUsuario1" ASC`, usuario.NombreUsuario)
	if err != nil {
		return []string{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var usuario string
		err = rows.Scan(&usuario)
		if err != nil {
			return []string{}, err
		}

		usuarios = append(usuarios, usuario)
	}
	return usuarios, err
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

// OtorgarPuntos añade una cantidad de puntos determinada al usuario dado. Devuelve error en caso de fallo.
func OtorgarPuntos(db *sql.DB, usuario *vo.Usuario, puntos int, partidaGanada bool) (err error) {
	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "puntos"="puntos"+$1 WHERE "nombreUsuario"=$2`, puntos, usuario.NombreUsuario)
	if err != nil {
		return err
	}

	err = AlmacenarNotificacionConEstado(db, usuario, logica_juego.NewNotificacionPuntosObtenidos(puntos, partidaGanada))
	// Un fallo de almacenar la notificación habiendo ya otorgado puntos no es crítico, se puede admitir continuar

	return err
}

// RetirarPuntos retira una cantidad de puntos determinada al usuario dado. Devuelve error en caso de fallo.
func RetirarPuntos(db *sql.DB, usuario *vo.Usuario, puntos int) (err error) {
	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "puntos"="puntos"-$1 WHERE "nombreUsuario"=$2`, puntos, usuario.NombreUsuario)

	return err
}

// ContabilizarPartidaGanada añade una partida ganada al usuario, contabilizándola también en el cómputo global
func ContabilizarPartidaGanada(db *sql.DB, usuario *vo.Usuario) (err error) {
	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "partidasGanadas"="partidasGanadas"+1 WHERE "nombreUsuario"=$1`, usuario.NombreUsuario)

	if err != nil {
		return err
	} else {
		return ContabilizarPartida(db, usuario)
	}
}

// ContabilizarPartida añade una partida jugada al usuario
func ContabilizarPartida(db *sql.DB, usuario *vo.Usuario) (err error) {
	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "partidasTotales"="partidasTotales"+1 WHERE "nombreUsuario"=$1`, usuario.NombreUsuario)

	return err
}

// ModificarBiografia actualiza la biografia del usuario
func ModificarBiografia(db *sql.DB, usuario *vo.Usuario, biografia string) error {
	_, err := db.Exec(`UPDATE backend."Usuario" SET biografia=$1 WHERE "nombreUsuario"=$2`, biografia, usuario.NombreUsuario)
	return err
}

// TieneObjeto devuelve true si y solo si el objeto "item" está en la colección de objetos de "usuario"
func TieneObjeto(db *sql.DB, usuario *vo.Usuario, item vo.ItemTienda) (existe bool, err error) {
	err = db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."TieneItems" WHERE "ID_item" = $1 AND "nombreUsuario" = $2)`,
		item.Id, usuario.NombreUsuario).Scan(&existe)
	if err != nil {
		return false, err
	}

	return existe, nil
}

// ModificarDados modifica el aspecto de dados equipado por el usuario
func ModificarDados(db *sql.DB, usuario *vo.Usuario, dados vo.ItemTienda) error {
	_, err := db.Exec(`UPDATE backend."Usuario" SET "ID_dado"=$1 WHERE "nombreUsuario"=$2`, dados.Id, usuario.NombreUsuario)
	return err
}

// ModificarFichas modifica el aspecto de fichas equipado por el usuario
func ModificarFichas(db *sql.DB, usuario *vo.Usuario, fichas vo.ItemTienda) error {
	_, err := db.Exec(`UPDATE backend."Usuario" SET "ID_ficha"=$1 WHERE "nombreUsuario"=$2`, fichas.Id, usuario.NombreUsuario)
	return err
}

// Ranking devuelve la lista de usuarios del sistema ordenada por partidas ganadas
func Ranking(db *sql.DB) (ranking []vo.ElementoRankingUsuarios, err error) {
	rows, err := db.Query(`SELECT "nombreUsuario", "partidasGanadas", "partidasTotales" FROM backend."Usuario" 
				ORDER BY "partidasGanadas" DESC`)

	if err != nil {
		return []vo.ElementoRankingUsuarios{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var jugador vo.ElementoRankingUsuarios
		err = rows.Scan(&jugador.NombreUsuario, &jugador.PartidasGanadas, &jugador.PartidasTotales)
		if err != nil {
			return []vo.ElementoRankingUsuarios{}, err
		}

		ranking = append(ranking, jugador)
	}

	return ranking, nil
}

// AlmacenarNotificacionConEstado guarda una notificación dependiente del estado del juego para el usuario dado.
// Se borrará junto al resto al ser consultadas en grupo
func AlmacenarNotificacionConEstado(db *sql.DB, usuario *vo.Usuario, notificacion interface{}) (err error) {
	err, notificaciones := ObtenerNotificacionesConEstado(db, usuario)
	if err != nil {
		return err
	}

	notificaciones = append(notificaciones, notificacion)

	err = almacenarNotificacionesConEstado(db, usuario, notificaciones)
	return err
}

// ObtenerNotificacionesConEstado devuelve un slice de notificaciones con estado almacenadas para el usuario.
// Todas las notificaciones se borrarán una vez consultadas
func ObtenerNotificacionesConEstado(db *sql.DB, usuario *vo.Usuario) (err error, notificaciones []interface{}) {
	var b bytes.Buffer
	decoder := gob.NewDecoder(&b)

	// Puede no tener notificaciones, y por tanto ser un NULL
	var buffer sql.NullString

	err = db.QueryRow(`SELECT "notificacionesPendientesConEstado" FROM backend."Usuario" WHERE "nombreUsuario"=$1`, usuario.NombreUsuario).Scan(&buffer)
	if err != nil {
		return err, notificaciones
	}

	if buffer.Valid { // No era NULL
		b.Write([]byte(buffer.String))
		err = decoder.Decode(&notificaciones)
	}

	// Borra las notificaciones con estado, ya que ya se han consultado
	_, err = db.Exec(`UPDATE backend."Usuario" SET "notificacionesPendientesConEstado" = NULL WHERE "nombreUsuario"=$1`, usuario.NombreUsuario)

	return err, notificaciones
}

// Función auxiliar para AlmacenarNotificacionConEstado
func almacenarNotificacionesConEstado(db *sql.DB, usuario *vo.Usuario, notificaciones []interface{}) (err error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err = encoder.Encode(notificaciones)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE "backend"."Usuario" SET "notificacionesPendientesConEstado" = $1 WHERE "backend"."Usuario"."nombreUsuario" = $2`, b.Bytes(), usuario.NombreUsuario)

	return err
}
