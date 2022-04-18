package handlers

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"net/http"
)

// TODO documentar
// Ruta: /api/consultarTienda
// Tipo: GET
func ConsultarTienda(writer http.ResponseWriter, request *http.Request) {
	objetos, err := dao.ConsultarTienda(globales.Db)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(objetos)
	escribirHeaderExito(writer)
}
