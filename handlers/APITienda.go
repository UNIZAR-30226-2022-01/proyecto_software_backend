package handlers

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/go-chi/chi/v5"
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

// TODO documentar
// consultar los objetos de un usuario
// Ruta: /api/consultarColeccion/{usuario}
// Tipo: GET
func ConsultarColeccion(writer http.ResponseWriter, request *http.Request) {
	usuario := chi.URLParam(request, "usuario")
	objetos, err := dao.ConsultarColeccion(globales.Db, usuario)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(objetos)
	escribirHeaderExito(writer)
}
