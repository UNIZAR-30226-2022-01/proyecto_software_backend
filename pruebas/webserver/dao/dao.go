package dao

import (
	"database/sql"
	_ "github.com/lib/pq" // Driver que usa el paquete de sql, para postgres
	"log"
)

// Constructor del objeto de base de datos
func InicializarConexionDb() *sql.DB {
	db, err := sql.Open("postgres", "postgres://golang:golang@postgres:5432/postgres?sslmode=disable")
	//db, err := sql.Open("postgres", "postgres://{user}:{password}@{hostname}:{port}/{database-name}?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	// Open hace un defer de abrir la conexión hasta que se intente ejecutar una query, por lo que se fuerza
	// a establecerla aquí por su hay algún error
	if err = db.Ping(); err != nil {
		log.Fatal("No se ha podido conectar a la BD:", err)
	}

	log.Println("Conectado a la DB.")

	return db
}
