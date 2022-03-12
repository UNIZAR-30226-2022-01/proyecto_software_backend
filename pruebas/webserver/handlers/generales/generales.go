// Package generales define handlers de páginas accesibles para cualquier usuario
package generales

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	//"html/template"
	//"log"
	//"net/http"
	//"webserver/globales"
	//"webserver/vo"
)

// Struct de ejemplo a serializar a JSON
type Ejemplo struct {
	Campo1 string
	Campo2 int
}

func HandlerDePrueba(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/index.html"))
	tmpl.Execute(writer, nil)
}

func HandlerDePruebaConParametros(writer http.ResponseWriter, request *http.Request) {
	month := chi.URLParam(request, "month")
	day := chi.URLParam(request, "day")
	year := chi.URLParam(request, "year")

	log.Println("Lo que estaba en la URL es", month, day, year)

	tmpl := template.Must(template.ParseFiles("html/index.html"))
	tmpl.Execute(writer, nil)
}

func MenuRegistro(writer http.ResponseWriter, request *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/formulario.html"))
	tmpl.Execute(writer, nil)
}

func HandlerDePruebaConParametrosPost(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("id")

	log.Println("Un parámetro del formulario es", id)

	tmpl := template.Must(template.ParseFiles("html/index.html"))
	tmpl.Execute(writer, nil)
}

func ServirJSON(writer http.ResponseWriter, request *http.Request) {
	ejemplo := Ejemplo{"campo1", 2}

	log.Println(json.Marshal(ejemplo))

	// Establece el contenido a servir como JSON y lo escribe
	json.NewEncoder(writer).Encode(ejemplo)
}

func ServirImagen(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "unizzard.png")
}
