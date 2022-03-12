// Package globales contiene variables globales a ser utilizadas por todos los módulos, instanciadas
// desde el paquete principal
package globales

import (
	"database/sql" // Funciones de sql
)

var Db *sql.DB // Base de datos thread safe, a compartir entre los módulos
