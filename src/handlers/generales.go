// Package generales define handlers de páginas accesibles para cualquier usuario y funciones auxiliares
package handlers

import (
	"log"
	"net/http"
)

// Devuelve una respuesta con status 500 junto al mensaje de error y la función
// en la que se ha dado.
func devolverError(writer http.ResponseWriter, funcion string, err error) {
	log.Println("Error en", funcion, ":", err)
	writer.WriteHeader(http.StatusInternalServerError)
	_, err = writer.Write([]byte(err.Error()))
	if err != nil {
		log.Println("Error al escribir respuesta en", funcion, ":", err)
	}
}

// Devuelve una respuesta con status 200.
func devolverExito(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
}
