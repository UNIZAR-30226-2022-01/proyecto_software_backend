package dao

import (
	"database/sql"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
)

// ConsultarTienda devuelve la lista de objetos disponibles en la tienda
func ConsultarTienda(db *sql.DB) (items []vo.ItemTienda, err error) {
	rows, err := db.Query(`SELECT id, nombre, descripcion, precio, tipo FROM backend."ItemTienda"`)
	if err != nil {
		return []vo.ItemTienda{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var item vo.ItemTienda
		err = rows.Scan(&item.Id, &item.Nombre, &item.Descripcion, &item.Precio, &item.Tipo)
		if err != nil {
			return []vo.ItemTienda{}, err
		}

		items = append(items, item)
	}

	return items, nil
}

// ConsultarColeccion permite consultar los objetos que ha comprado un usuario
func ConsultarColeccion(db *sql.DB, usuario string) (items []vo.ItemTienda, err error) {
	rows, err := db.Query(`SELECT id, nombre, descripcion, precio, tipo FROM backend."ItemTienda" 
    	JOIN backend."TieneItems" TI on "ItemTienda".id = TI."ID_item" and TI."nombreUsuario" = $1`, usuario)
	if err != nil {
		return []vo.ItemTienda{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var item vo.ItemTienda
		err = rows.Scan(&item.Id, &item.Nombre, &item.Descripcion, &item.Precio, &item.Tipo)
		if err != nil {
			return []vo.ItemTienda{}, err
		}

		items = append(items, item)
	}

	return items, nil
}

// ObtenerIDAvatar devuelve el ID de item del avatar equipado para el usuario dado
func ObtenerIDAvatar(db *sql.DB, usuario string) (id int, err error) {
	err = db.QueryRow(`SELECT "ID_avatar" FROM backend."Usuario" WHERE "nombreUsuario" = $1`, usuario).Scan(&id)

	return id, err
}

// ObtenerIDDado devuelve el ID de item de los dados equipados para el usuario dado
func ObtenerIDDado(db *sql.DB, usuario string) (id int, err error) {
	err = db.QueryRow(`SELECT "ID_dado" FROM backend."Usuario" WHERE "nombreUsuario" = $1`, usuario).Scan(&id)

	return id, err
}

// ObtenerObjeto recupera un objeto de la tienda de la base de datos a partir de su identificador
func ObtenerObjeto(db *sql.DB, idItem int) (vo.ItemTienda, error) {
	var item vo.ItemTienda
	err := db.QueryRow(`SELECT id, nombre, descripcion, precio, tipo FROM backend."ItemTienda" WHERE id = $1`,
		idItem).Scan(&item.Id, &item.Nombre, &item.Descripcion, &item.Precio, &item.Tipo)

	if err != nil {
		return vo.ItemTienda{}, err
	}

	return item, nil
}

// ComprarObjeto permite al jugador comprar un objeto de la tienda siempre y cuando tenga los puntos necesarios.
// Para ello, se especificar?? como par??metro el identificador del objeto que desea comprar. La compra se realizar??
// siempre que dicho objeto exista, no sea uno de los objetos iniciales, el jugador tenga los puntos suficientes para
// comprarlo y el jugador no lo haya comprado ya.
// Si se indica que se est?? creando el usuario, se permite comprar objetos iniciales
func ComprarObjeto(db *sql.DB, usuario string, item vo.ItemTienda, creacion bool) error {
	// Comprobar que el objeto existe
	var existe bool
	err := db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."ItemTienda" WHERE id = $1)`, item.Id).Scan(&existe)
	if err != nil {
		return err
	}

	if !existe {
		return errors.New("El objeto no existe")
	}

	// Comprobar que no es un objeto inicial
	if item.Precio == 0 && !creacion {
		return errors.New("No puedes comprar uno de los objetos iniciales")
	}

	// Comprobar que no tiene el objeto
	err = db.QueryRow(`SELECT EXISTS(SELECT * FROM backend."TieneItems" WHERE "ID_item" = $1 AND "nombreUsuario" = $2)`,
		item.Id, usuario).Scan(&existe)
	if err != nil {
		return err
	}

	if existe {
		return errors.New("No puedes comprar un objeto que ya tienes")
	}

	// Comprobar que tiene los puntos suficientes
	var puntosUsuario int
	err = db.QueryRow(`SELECT puntos FROM backend."Usuario" WHERE "nombreUsuario" = $1`, usuario).Scan(&puntosUsuario)
	if err != nil {
		return err
	}

	if item.Precio > puntosUsuario {
		return errors.New("No tienes puntos suficientes para comprar el objeto")
	}

	// Comprar el objeto
	err = RetirarPuntos(db, &vo.Usuario{NombreUsuario: usuario}, item.Precio)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO backend."TieneItems" VALUES ($1, $2)`, item.Id, usuario)
	if err != nil {
		return err
	}

	return nil
}
