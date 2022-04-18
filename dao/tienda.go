package dao

import (
	"database/sql"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
)

// TODO documentar
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

// TODO documentar
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
