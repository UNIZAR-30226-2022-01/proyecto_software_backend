// Package generales define handlers de páginas accesibles para cualquier usuario y funciones auxiliares
package handlers

import (
	"errors"
	"log"
	"net/http"
)

// Devuelve una respuesta con status 500 junto al mensaje de error y la función
// en la que se ha dado.
func devolverError(writer http.ResponseWriter, err error) {
	log.Println("Error:", err)
	writer.WriteHeader(http.StatusInternalServerError)
	_, err = writer.Write([]byte(err.Error()))
	if err != nil {
		log.Println("Error al escribir respuesta en:", err)
	}
}

// Devuelve una respuesta con status 500 junto al mensaje de error y la función
// en la que se ha dado.
func devolverErrorSQL(writer http.ResponseWriter) {
	err := errors.New("Se ha producido un error en la base de datos.")

	log.Println("Error en:", err)
	writer.WriteHeader(http.StatusInternalServerError)
	_, err = writer.Write([]byte(err.Error()))
	if err != nil {
		log.Println("Error al escribir respuesta en:", err)
	}
}

// Devuelve una respuesta con status 200.
func devolverExito(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
}
