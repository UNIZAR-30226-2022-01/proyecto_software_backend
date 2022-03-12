// Package middleware define middleware propio para actuar de intermediario entre
// la llegada de una petición y su tratamiento final por un handler
package middleware

import (
	"log"
	"net/http"
)

// Función que devuelve una función de middleware
func MiddlewarePropio() func(next http.Handler) http.Handler {
	// next es el handler (o middleware) siguiente a éste middleware
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Escribe directamente y luego deja escribir al handler a continuación,
			// en realidad se leerían cookies y se serviría contenido diferente o
			// dejaría pasar al handler, etc.
			//w.Write([]byte("Hola desde el middleware!"))

			log.Println("Hola desde el middleware!")

			// Deja pasar al handler
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
