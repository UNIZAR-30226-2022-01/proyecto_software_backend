// Package dao define funciones de comunicación directa con la base de datos
package dao

import (
	"database/sql"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	_ "github.com/lib/pq" // Driver que usa el paquete de sql, para postgres
	"log"
	"os"
	"time"
)

// InicializarConexionDb devuelve el objeto de base de datos, en el cual realiza la conexión a la misma
func InicializarConexionDb(test bool) *sql.DB {
	//db, err := sql.Open("postgres", "postgres://{user}:{password}@{hostname}:{port}/{database-name}?sslmode=disable")
	var db *sql.DB
	var err error
	if test {
		// Para pruebas fuera de Docker:
		db, err = sql.Open("postgres", "postgres://"+os.Getenv(globales.USUARIO_DB)+":"+os.Getenv(globales.PASSWORD_DB)+"@"+os.Getenv(globales.DIRECCION_DB_TESTS)+":5432/postgres?sslmode=disable")
	} else {
		db, err = sql.Open("postgres", "postgres://"+os.Getenv(globales.USUARIO_DB)+":"+os.Getenv(globales.PASSWORD_DB)+"@"+os.Getenv(globales.DIRECCION_DB)+":5432/postgres?sslmode=disable")
	}

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

// MonitorizarCanalBorrado pone en marcha una Goroutine de atención a eliminado de partidas y usuarios. Diseñado
// para eliminar partidas terminadas o con usuarios inactivos y usuarios inactivos, sin necesitar la intervención de los handlers.
func MonitorizarCanalBorrado(db *sql.DB, partidas chan int, stop chan struct{}, usuariosInactivos chan string) {
	var idPartida int
	var jugador string

	go func(db *sql.DB, partidas chan int, stop chan struct{}, usuariosInactivos chan string) {
		for {
			// Si hubiera un fallo en la base de datos que impidiera el borrado de una partida o usuario,
			// tras el reinicio del servidor se reintentaría, ya que se cargará de nuevo las partidas en la cache
			select {
			case idPartida = <-partidas:
				go BorrarPartida(db, &vo.Partida{IdPartida: idPartida})
			case jugador = <-usuariosInactivos:
				go AbandonarPartida(db, &vo.Usuario{NombreUsuario: jugador})
				go AlmacenarNotificacionConEstado(db, &vo.Usuario{NombreUsuario: jugador}, logica_juego.NewNotificacionExpulsion())
			case <-stop:
				return
			}
		}
	}(db, partidas, stop, usuariosInactivos)
}
