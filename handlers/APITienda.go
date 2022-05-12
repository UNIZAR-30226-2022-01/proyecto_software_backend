package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// ConsultarTienda devuelve la lista de objetos disponibles en la tienda en formato JSON. La estructura de la lista
// devuelta es la siguiente:
// [ {
//    Id: 			int
// 	  Nombre: 		string
// 	  Descripcion: 	string
// 	  Precio: 		int
// 	  Tipo: 		string ({"avatar", "dado"})
//	  Imagen:		string (blob)
//	  }, {...}, ...
// ]
//
// Devuelve status 500 en caso de error y 200 en caso contrario
//
// Ruta: /api/consultarTienda
// Tipo: GET
func ConsultarTienda(writer http.ResponseWriter, request *http.Request) {
	objetos, err := dao.ConsultarTienda(globales.Db)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Se cargan las imágenes de los objetos
	var bytes []byte
	for i, objeto := range objetos {
		if objeto.Tipo == globales.TIPO_AVATAR {
			bytes, err = ioutil.ReadFile(globales.RUTA_AVATARES + strconv.Itoa(objeto.Id) + globales.FORMATO_ASSETS)
			if err != nil {
				log.Println("error al cargar img:", err)
				devolverErrorSQL(writer)
				return
			}
		} else if objeto.Tipo == globales.TIPO_DADO {
			bytes, err = ioutil.ReadFile(globales.RUTA_DADOS + strconv.Itoa(objeto.Id) + "5" + globales.FORMATO_ASSETS)
			if err != nil {
				log.Println("error al cargar img:", err)
				devolverErrorSQL(writer)
				return
			}
		}

		objetos[i].Imagen = bytes
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(objetos)
	escribirHeaderExito(writer)
}

// ConsultarColeccion permite consultar los objetos que ha comprado un usuario, cuyo nombre se especifica como parte de
// la URL. Devuelve una lista de objetos codificada en JSON, con el siguiente formato:
// [ {
//    Id: 			int
// 	  Nombre: 		string
// 	  Descripcion: 	string
// 	  Precio: 		int
// 	  Tipo: 		string ({"avatar", "dado"})
//	  Imagen:		string (blob)
//	  }, {...}, ...
// ]
//
// Devuelve status 500 en caso de error y 200 en caso contrario
//
// Ruta: /api/consultarColeccion/{usuario}
// Tipo: GET
func ConsultarColeccion(writer http.ResponseWriter, request *http.Request) {
	usuario := chi.URLParam(request, "usuario")
	objetos, err := dao.ConsultarColeccion(globales.Db, usuario)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	// Se cargan las imágenes de los objetos
	var bytes []byte
	for i, objeto := range objetos {
		if objeto.Tipo == globales.TIPO_AVATAR {
			bytes, err = ioutil.ReadFile(globales.RUTA_AVATARES + strconv.Itoa(objeto.Id) + globales.FORMATO_ASSETS)
			if err != nil {
				log.Println("error al cargar img:", err)
				devolverErrorSQL(writer)
				return
			}
		} else if objeto.Tipo == globales.TIPO_DADO {
			bytes, err = ioutil.ReadFile(globales.RUTA_DADOS + strconv.Itoa(objeto.Id) + "5" + globales.FORMATO_ASSETS)
			if err != nil {
				log.Println("error al cargar img:", err)
				devolverErrorSQL(writer)
				return
			}
		}

		objetos[i].Imagen = bytes
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(objetos)
	escribirHeaderExito(writer)
}

// ObtenerAvatar devuelve una imagen codificada como octet-stream para el avatar
// del usuario indicado.
//
// Devuelve un error 500 en caso de que el usuario no exista u ocurra cualquier otro error
//
// Ruta: /api/obtenerFotoPerfil/{usuario}
// Tipo: GET
func ObtenerAvatar(writer http.ResponseWriter, request *http.Request) {
	usuario := chi.URLParam(request, "usuario")
	idAvatar, err := dao.ObtenerIDAvatar(globales.Db, usuario)
	if err != nil {
		log.Println("error al obtener id:", err)
		devolverErrorSQL(writer)
		return
	}

	bytes, err := ioutil.ReadFile(globales.RUTA_AVATARES + strconv.Itoa(idAvatar) + globales.FORMATO_ASSETS)
	if err != nil {
		log.Println("error al cargar img:", err)
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/octet-stream")
	_, err = writer.Write(bytes)
	if err != nil {
		log.Println("error al escribir img:", err)
		devolverErrorSQL(writer)
	} else {
		escribirHeaderExito(writer)
	}
}

// ObtenerDados devuelve una imagen codificada como octet-stream para la cara indicada de los dados
// equipados del usuario indicado. La cara debe ser un valor entre 1 y 6, correspondiente
// al valor de los dados.
//
// Devuelve un error 500 en caso de que la cara sea inválida, el usuario no exista u ocurra cualquier otro error.
//
// Ruta: /api/obtenerDados/{usuario}/{cara}
// Tipo: GET
func ObtenerDados(writer http.ResponseWriter, request *http.Request) {
	usuario := chi.URLParam(request, "usuario")
	cara := chi.URLParam(request, "cara")
	idDados, err := dao.ObtenerIDDado(globales.Db, usuario)

	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	bytes, err := ioutil.ReadFile(globales.RUTA_DADOS + strconv.Itoa(idDados) + cara + globales.FORMATO_ASSETS)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	writer.Header().Set("Content-Type", "application/octet-stream")
	_, err = writer.Write(bytes)
	if err != nil {
		devolverErrorSQL(writer)
	} else {
		escribirHeaderExito(writer)
	}
}

// ObtenerImagenItem devuelve la imagen de un ítem, dado su ID, codificado como blob (octet-stream)
//
// Devuelve un error 500 en caso de que el ID sea inválido u ocurra cualquier otro error
//
// Ruta: /api/obtenerImagenItem/{id}
// Tipo: GET
func ObtenerImagenItem(writer http.ResponseWriter, request *http.Request) {
	var bytes []byte
	idParam := chi.URLParam(request, "id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		devolverError(writer, errors.New("El id de ítem debe ser un número entero"))
		return
	}

	item, err := dao.ObtenerObjeto(globales.Db, id)
	if err != nil {
		devolverError(writer, errors.New("El id de ítem no existe"))
		return
	}

	if item.Tipo == globales.TIPO_AVATAR {
		bytes, err = ioutil.ReadFile(globales.RUTA_AVATARES + strconv.Itoa(item.Id) + globales.FORMATO_ASSETS)
		if err != nil {
			log.Println("error al cargar img:", err)
			devolverErrorSQL(writer)
			return
		}
	} else if item.Tipo == globales.TIPO_DADO {
		bytes, err = ioutil.ReadFile(globales.RUTA_DADOS + strconv.Itoa(item.Id) + "5" + globales.FORMATO_ASSETS)
		if err != nil {
			log.Println("error al cargar img:", err)
			devolverErrorSQL(writer)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/octet-stream")
	_, err = writer.Write(bytes)
}

// ComprarObjeto permite al jugador comprar un objeto de la tienda siempre y cuando tenga los puntos necesarios.
// Para ello, especificará como parte de la URL el identificador del objeto que desea comprar. La compra se realizará
// siempre que dicho objeto exista, no sea uno de los objetos iniciales, el jugador tenga los puntos suficientes para
// comprarlo y el jugador no lo haya comprado ya.
//
// Devuelve status 500 en caso de error y 200 en caso contrario
//
// Ruta: /api/comprarObjeto/{id_objeto}
// Tipo: POST
func ComprarObjeto(writer http.ResponseWriter, request *http.Request) {
	idItem, err := strconv.Atoi(chi.URLParam(request, "id_objeto"))
	usuario := middleware.ObtenerUsuarioCookie(request)

	if err != nil {
		devolverError(writer, errors.New("El identificador del objeto debe ser un número natural"))
		return
	}

	item, err := dao.ObtenerObjeto(globales.Db, idItem)
	if err != nil {
		devolverErrorSQL(writer)
		return
	}

	err = dao.ComprarObjeto(globales.Db, usuario, item, false)
	if err == sql.ErrNoRows {
		devolverErrorSQL(writer)
		return
	}
	if err != nil {
		devolverError(writer, err)
		return
	}

	escribirHeaderExito(writer)
}
