package dao

import (
	"database/sql"
	_ "github.com/lib/pq" // Driver que usa el paquete de sql, para postgres
	"log"
	"time"
)

// InicializarConexionDb devuelve el objeto de base de datos, en el cual realiza la conexión a la misma
func InicializarConexionDb() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable")
	//db, err := sql.Open("postgres", "postgres://{user}:{password}@{hostname}:{port}/{database-name}?sslmode=disable")

	// Para pruebas fuera de Docker:
	//db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	// Open hace un defer de abrir la conexión hasta que se intente ejecutar una query, por lo que se fuerza
	// a establecerla aquí por su hay algún error
	if err = db.Ping(); err != nil {
		// Reintenta si la primera conexión no tiene éxito, posiblemente debido a que se ha adelantado al contenedor de
		// postgres en el intervalo en el que está en marcha pero aún no atiende peticiones
		time.Sleep(5 * time.Second)

		if err = db.Ping(); err != nil {
			log.Fatal("No se ha podido conectar a la BD:", err)
		}
	}

	log.Println("Conectado a la DB.")

	return db
}
