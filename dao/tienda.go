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
		err = rows.Scan(&item)
		if err != nil {
			return []vo.ItemTienda{}, err
		}

		items = append(items, item)
	}

	return items, nil
}
