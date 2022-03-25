// Package globales contiene variables globales a ser utilizadas por todos los módulos, instanciadas
// desde el paquete principal
package globales

import (
	"database/sql" // Funciones de sql
)

var Db *sql.DB // Base de datos thread safe, a compartir entre los módulos

const (
	DIRECCION_DB       = "DIRECCION_DB"
	DIRECCION_DB_TESTS = "DIRECCION_DB_TESTS"
	PUERTO_WEB         = "PUERTO_WEB"
	PUERTO_API         = "PUERTO_API"
	USUARIO_DB         = "USUARIO_DB"
	PASSWORD_DB        = "PASSWORD_DB"
	CARPETA_FRONTEND   = "web"
)
