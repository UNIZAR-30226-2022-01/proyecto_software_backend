package main

import (
	"backend/dao"
	"backend/globales"
	"backend/handlers"
	middlewarePropio "backend/middleware"
	"context"
	"io/ioutil"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	/////////////////////////////////
	// Librerías externas
	////////////////////////////////
	// Logging de errores y acciones
	"log"
	// Módulo de HTTP estándar
	"net/http"

	// Librería auxiliar a la estándar que hace expande las capacidades del router de peticiones
	// HTTP, permitiendo establecer diferentes respuestas para POST/GET/DELETE, usar regex en URLs,
	// establecer middlewares para grupos o URLs individuales, obtener más fácilmente parámetros de
	// URLs pre-establecidos, etc.
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"os"
	// Middleware a utilizar escrito por nosotros
)

const (
	CARPETA_FRONTEND      = "web"
	FICHERO_RAIZ_FRONTEND = "index.html"
)

func main() {
	var server *http.Server
	if len(os.Args) < 2 {
		log.Println("Uso:\n\t ./ejecutable -web : Servir contenido web \n\t ./ejecutable -api : Servir API")
		os.Exit(1)
	} else {
		// Instancia un servidor HTTP con el router programado indicado
		if os.Args[1] == "-web" {
			server = &http.Server{Addr: ":8080", Handler: routerWeb()}
		} else if os.Args[1] == "-api" {
			server = &http.Server{Addr: ":8080", Handler: routerAPI()}

			// El objeto de base de datos es seguro para uso concurrente y controla su
			// propia pool de conexiones independientemente.
			globales.Db = dao.InicializarConexionDb()
		} else {
			log.Println("Uso:\n\t ./ejecutable -web : Servir contenido web \n\t ./ejecutable -api : Servir API")
			os.Exit(1)
		}
	}

	canalCierre := tratarContextoCierreServidor(server)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Espera a que el servidor esté cerrado
	<-canalCierre

	// Termina todos los módulos de forma segura
	if os.Args[1] == "-api" {
		globales.Db.Close()
	}

	os.Exit(0)
}

// Crea un contexto de cierre de un servidor HTTP y una goroutina de trata del mismo, y
// devuelve un canal de espera para su cierre
func tratarContextoCierreServidor(server *http.Server) <-chan struct{} {
	// Crea el contexto y la función de cancelación para el servidor
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Crea un canal y una función de tratamiento de señales de cancelación
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Crea un contexto para indicar al servidor de que debería terminar antes de 30 segundos
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Se ha agotado el tiempo de gracia para el cierre del servidor. Terminando el proceso forzosamente...")
			}
		}()

		// Indica al servidor que deje de atender y termine las peticiones en curso
		// con el tiempo límite del contexto dado
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}

		// Una vez parado, se cierra el servidor
		serverStopCtx()
	}()

	return serverCtx.Done()
}

// Devuelve un router programado para las URLs a atender
func routerAPI() http.Handler {
	r := chi.NewRouter()

	// Para debugging
	r.Use(middleware.Logger)

	// Formularios
	r.Post("/registro", handlers.Registro)
	r.Post("/login", handlers.Login)
	//TODO: Otro POST para formularios de cambiar perfil de usuario

	// Pruebas
	r.Get("/formularioRegistro", handlers.MenuRegistro)
	r.Get("/formulariologin", handlers.MenuLogin)

	// Rutas REST
	r.Route("/api", func(r chi.Router) {
		// Obligamos el acceso con login previo
		r.Use(middlewarePropio.MiddlewareSesion())

		// Partidas
		r.Post("/crearPartida", handlers.CrearPartida)
		r.Post("/unirseAPartida", handlers.UnirseAPartida) // TODO: Mejor con URL?
		r.Get("/obtenerPartidas", handlers.ObtenerPartidas)

		// Usuarios
		r.Post("/aceptarSolicitudAmistad/{nombre}", handlers.AceptarSolicitudAmistad)
		r.Post("/rechazarSolicitudAmistad/{nombre}", handlers.RechazarSolicitudAmistad)
		r.Post("/enviarSolicitudAmistad/{nombre}", handlers.EnviarSolicitudAmistad)
		r.Get("/obtenerNotificaciones/", handlers.ObtenerNotificaciones)
	})

	return r
}

func routerWeb() http.Handler {
	r := chi.NewRouter()

	// Para debugging
	r.Use(middleware.Logger)

	directorioDeTrabajo, _ := os.Getwd()
	ficherosFrontend := filepath.Join(directorioDeTrabajo, CARPETA_FRONTEND)
	log.Println("Sirviendo " + FICHERO_RAIZ_FRONTEND + " desde " + ficherosFrontend)
	index, _ := ioutil.ReadFile(ficherosFrontend + "/" + FICHERO_RAIZ_FRONTEND)
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	})

	return r
}
