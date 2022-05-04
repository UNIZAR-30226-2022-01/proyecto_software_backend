// Package dao define funciones de comunicación directa con la base de datos
package dao

import (
	"crypto/tls"
	"database/sql"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	_ "github.com/lib/pq" // Driver que usa el paquete de sql, para postgres
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	LONGITUD_TOKEN_RESET_PASSWORD = 40
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

// MonitorizarCanalEnvioAlertas pone en marcha una GOroutine de atención a envío de emails a usuarios. Diseñado para ser
// utilizado por el gestor de cache de partidas, evitando dependencias cíclicas en el dao.
func MonitorizarCanalEnvioAlertas(db *sql.DB, stop chan struct{}, jugadores chan string) {
	var jugador string

	go func(db *sql.DB, stop chan struct{}, jugadores chan string) {
		for {
			// Si hubiera un fallo en la base de datos que impidiera el borrado de una partida o usuario,
			// tras el reinicio del servidor se reintentaría, ya que se cargará de nuevo las partidas en la cache
			select {
			case jugador = <-jugadores:
				go enviarAlerta(jugador)
			case <-stop:
				return
			}
		}
	}(db, stop, jugadores)
}

func enviarAlerta(jugador string) {
	_, email := ObtenerEmailUsuario(globales.Db, jugador)
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv(globales.DIRECCION_ENVIO_EMAILS))
	m.SetHeader("To", email)
	m.SetHeader("Subject", jugador+", tienes una partida pendiente de que hagas tu movimiento.")
	//m.SetBody("text/html", "<a>"+os.Getenv(globales.NOMBRE_DNS_API)+"/resetearPassword/</a>"+token)
	m.SetBody("text/html", "¡Accede a <a>"+os.Getenv(globales.NOMBRE_DNS_REACT)+"</a> o <a>"+os.Getenv(globales.NOMBRE_DNS_ANGULAR)+"</a> para continuar tu partida!")

	puerto, _ := strconv.Atoi(os.Getenv(globales.PUERTO_SMTP))
	d := gomail.NewDialer(os.Getenv(globales.HOST_SMTP), puerto, os.Getenv(globales.USUARIO_SMTP), os.Getenv(globales.PASS_SMTP))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Evita problemas derivados de no tener certificados en el contenedor de Docker

	// Como no se puede verificar que el destino no existe, los errores al enviar correos se ignoran silenciosamente
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error al enviar email de reset de contraseña a", jugador, ":", err)
	}
}
