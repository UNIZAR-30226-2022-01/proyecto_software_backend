package servidor

import (
	"backend/dao"
	"backend/globales"
	"backend/handlers"
	"backend/logica_juego"
	"backend/vo"

	middlewarePropio "backend/middleware" // Middleware a utilizar escrito por nosotros
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func IniciarServidor(test bool) {
	var server *http.Server
	if len(os.Args) < 2 && !test {
		log.Println("Uso:\n\t ./ejecutable -web : Servir contenido web \n\t ./ejecutable -api : Servir API")
		os.Exit(1)
	} else {
		// Instancia un servidor HTTP con el router programado indicado
		if os.Args[len(os.Args)-1] == "-web" {
			log.Println("Escuchando por el puerto", os.Getenv(globales.PUERTO_WEB))
			server = &http.Server{Addr: ":" + os.Getenv(globales.PUERTO_WEB), Handler: routerWeb()}
		} else if os.Args[len(os.Args)-1] == "-api" || test {
			log.Println("Escuchando por el puerto", os.Getenv(globales.PUERTO_API))
			server = &http.Server{Addr: ":" + os.Getenv(globales.PUERTO_API), Handler: routerAPI()}

			// El objeto de base de datos es seguro para uso concurrente y controla su
			// propio pool de conexiones independientemente.
			globales.Db = dao.InicializarConexionDb(test)
		} else {
			log.Println("Uso:\n\t ./ejecutable -web : Servir contenido web \n\t ./ejecutable -api : Servir API")
			os.Exit(1)
		}
	}

	canalCierre := tratarContextoCierreServidor(server)

	// Inicio de lógica del juego
	globales.InicializarGrafoMapa()
	globales.AlmacenPartidas = logica_juego.IniciarAlmacenPartidas()

	go func(cs chan vo.Partida, cp chan struct{}) {
		for {
			select {
			case partida := <-cs:
				dao.AlmacenarEstadoSerializado(globales.Db, &partida)
			case <-cp:
				break
			}
		}
	}(globales.AlmacenPartidas.CanalSerializacion, globales.AlmacenPartidas.CanalParada)

	partidas, err := dao.ObtenerPartidas(globales.Db)
	if err != nil {
		log.Fatal("Error al recuperar partidas almacenadas:", err)
	}

	for _, p := range partidas {
		globales.AlmacenPartidas.AlmacenarPartida(p)
	}

	err = server.ListenAndServe()
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

// Devuelve un router programado para URLs de la API
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
		r.Post("/unirseAPartida", handlers.UnirseAPartida)
		r.Post("/abandonarLobby", handlers.AbandonarLobby)
		r.Get("/obtenerPartidas", handlers.ObtenerPartidas)

		// Usuarios
		r.Post("/aceptarSolicitudAmistad/{nombre}", handlers.AceptarSolicitudAmistad)
		r.Post("/rechazarSolicitudAmistad/{nombre}", handlers.RechazarSolicitudAmistad)
		r.Post("/enviarSolicitudAmistad/{nombre}", handlers.EnviarSolicitudAmistad)
		r.Get("/obtenerNotificaciones/", handlers.ObtenerNotificaciones)
	})

	return r
}

// Devuelve un router programado para URLs de cualquiera de los frontends
func routerWeb() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	workDir, _ := os.Getwd()
	// Carpeta del sistema de ficheros que se va a servir, restringida a ella y sus
	// subdirectorios
	filesDir := http.Dir(filepath.Join(workDir, globales.CARPETA_FRONTEND))
	fileServer(r, "/", filesDir)

	return r
}

// fileServer pone en marcha un router para un servidor de ficheros mediante HTTP,
// que sirve ficheros estáticos desde root
// Adaptado de https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("No se permite servir desde directorios con parámetros de URL.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
