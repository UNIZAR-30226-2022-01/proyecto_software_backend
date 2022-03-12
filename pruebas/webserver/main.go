package main

import (
	"context"
	"os/signal"
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
	/////////////////////////////////
	// Módulos propios
	/////////////////////////////////
	// Funciones de tratamiento de peticiones HTTP
	"backend/handlers/generales"
	// Middleware a utilizar escrito por nosotros
	middlewarePropio "backend/middleware"
)

func main() {
	// El objeto de base de datos es seguro para uso concurrente y controla su
	// propia pool de conexiones independientemente.
	//globales.Db = dao.InicializarConexionDb()

	// Instancia un servidor HTTP con el router programado indicado
	server := &http.Server{Addr: ":8080", Handler: service()}

	canalCierre := tratarContextoCierreServidor(server)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Espera a que el servidor esté cerrado
	<-canalCierre

	// Termina todos los módulos de forma segura
	//globales.Db.Close()
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
func service() http.Handler {
	r := chi.NewRouter()

	// Middlewares a usar, de momento sobre todas las URLs
	r.Use(middleware.Logger)
	r.Use(middlewarePropio.MiddlewarePropio())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenido!"))
	})

	// Diferentes pruebas de las funcionalidaes de Chi
	r.Post("/registro", generales.HandlerDePruebaConParametrosPost)
	r.Get("/dameJSON", generales.ServirJSON)
	r.Get("/dameUnaImagen", generales.ServirImagen)
	r.Get("/formularioRegistro", generales.MenuRegistro)

	// Ejemplo de rutas REST
	r.Route("/articles", func(r chi.Router) {
		r.Get("/", generales.HandlerDePrueba)                                  // GET /articles
		r.Get("/{month}-{day}-{year}", generales.HandlerDePruebaConParametros) // GET /articles/01-16-2017

	})

	return r
}
