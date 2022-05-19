// Package servidor define y controla el servidor HTTP y las rutas a la API o al contenido web, según
// los argumentos introducidos
package servidor

import (
	"context"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/handlers"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	middlewarePropio "github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware" // Middleware a utilizar escrito por nosotros
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func IniciarServidor(test bool) {
	// Inicializa el gestor de certificados
	// Tiene como dominios permitidos el global y los subdominios de la API, angular y react, y pedirá los certificados a
	// Let's Encrypt dinámicamente dependiendo de la petición
	//
	// Los certificados se cachearán en el directorio indicado
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(os.Getenv(globales.NOMBRE_DNS_GLOBAL), os.Getenv(globales.NOMBRE_DNS_API), os.Getenv(globales.NOMBRE_DNS_ANGULAR), os.Getenv(globales.NOMBRE_DNS_REACT)),
		Cache:      autocert.DirCache(globales.RUTA_CACHE_CERTIFICADOS),
	}

	log.Println("Configurado certificado TLS para:")
	log.Println(os.Getenv(globales.NOMBRE_DNS_GLOBAL))
	log.Println(os.Getenv(globales.NOMBRE_DNS_API))
	log.Println(os.Getenv(globales.NOMBRE_DNS_ANGULAR))
	log.Println(os.Getenv(globales.NOMBRE_DNS_REACT))

	tlsConfig := certManager.TLSConfig()

	var server *http.Server
	if len(os.Args) < 2 && !test {
		log.Println("Uso:\n\t ./ejecutable -web : Servir contenido web \n\t ./ejecutable -api : Servir API")
		os.Exit(1)
	} else {
		// Instancia un servidor HTTP con el router programado indicado
		if os.Args[len(os.Args)-1] == "-web" {
			log.Println("Escuchando por el puerto", os.Getenv(globales.PUERTO_WEB))
			server = &http.Server{Addr: ":" + os.Getenv(globales.PUERTO_WEB),
				Handler:   routerWeb(),
				TLSConfig: tlsConfig}
		} else if os.Args[len(os.Args)-1] == "-api" || test {
			log.Println("Escuchando por el puerto", os.Getenv(globales.PUERTO_API))
			server = &http.Server{Addr: ":" + os.Getenv(globales.PUERTO_API),
				Handler:   routerAPI(),
				TLSConfig: tlsConfig}

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
	if os.Args[len(os.Args)-1] == "-api" || test {
		logica_juego.InicializarGrafoMapa()
		logica_juego.InicializarContinentes()
		globales.CachePartidas = globales.IniciarAlmacenPartidas()
		globales.IniciarCanalesEliminacionPartidasDB()
		dao.MonitorizarCanalBorrado(globales.Db, globales.CanalEliminacionPartidasDB, globales.CanalParadaBorradoPartidasDB, globales.CanalExpulsionUsuariosDB)
		dao.MonitorizarCanalEnvioAlertas(globales.Db, globales.CanalParadaEnvioAlertas, globales.CanalEnvioAlertas)
		// Registra los tipos a decodificar por gob a partir de interface{}
		logica_juego.RegistrarAcciones()
		logica_juego.RegistrarNotificaciones()

		go func(cs chan vo.Partida, cp chan struct{}) {
			for {
				select {
				case partida := <-cs:
					err := dao.AlmacenarEstadoSerializado(globales.Db, &partida)

					if err != nil { // Se ha roto la consistencia, no se puede seguir
						log.Fatal("Error al almacenar estado serializado:", err)
					}
				case <-cp:
					break
				}
			}
		}(globales.CachePartidas.CanalSerializacion, globales.CachePartidas.CanalParada)

		partidas, err := dao.ObtenerPartidas(globales.Db)
		if err != nil {
			log.Fatal("Error al recuperar partidas almacenadas:", err)
		}

		for _, p := range partidas {
			globales.CachePartidas.AlmacenarPartida(p)
		}
	}

	err := server.ListenAndServeTLS(globales.RUTA_CERT_TLS, globales.RUTA_CLAVE_TLS)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Espera a que el servidor esté cerrado
	<-canalCierre

	// Termina todos los módulos de forma segura
	if os.Args[1] == "-api" {
		globales.Db.Close()
		globales.CanalParadaBorradoPartidasDB <- struct{}{}
		globales.CanalParadaEnvioAlertas <- struct{}{}
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST"}, //"PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Formularios
	r.Post("/registro", handlers.Registro)
	r.Post("/login", handlers.Login)
	r.Post("/resetearPassword", handlers.ResetearContraseña)
	r.Post("/obtenerTokenResetPassword", handlers.ObtenerTokenResetPassword)

	// Rutas REST
	r.Route("/api", func(r chi.Router) {
		// Obligamos el acceso con login previo
		r.Use(middlewarePropio.MiddlewareSesion())

		// Partidas y lobby
		r.Post("/crearPartida", handlers.CrearPartida)
		r.Post("/unirseAPartida", handlers.UnirseAPartida)
		r.Post("/abandonarLobby", handlers.AbandonarLobby)
		r.Post("/abandonarPartida", handlers.AbandonarPartida)
		r.Get("/obtenerPartidas", handlers.ObtenerPartidas)
		r.Get("/obtenerEstadoPartida", handlers.ObtenerEstadoPartida)
		r.Get("/obtenerEstadoPartidaCompleto", handlers.ObtenerEstadoPartidaCompleto)
		r.Get("/resumirPartida", handlers.ResumirPartida)
		r.Get("/jugandoEnPartida", handlers.JugandoEnPartida)
		r.Get("/obtenerJugadoresPartida", handlers.ObtenerJugadoresPartida)
		r.Post("/enviarMensaje", handlers.EnviarMensaje)

		// Acciones del juego
		r.Get("/obtenerEstadoLobby", handlers.ObtenerEstadoLobby)
		r.Get("/consultarCartas", handlers.ConsultarCartas)
		r.Post("/reforzarTerritorio/{id}/{numTropas}", handlers.ReforzarTerritorio)
		r.Post("/cambiarCartas/{carta1}/{carta2}/{carta3}", handlers.CambiarCartas)
		r.Post("/pasarDeFase", handlers.PasarDeFase)
		r.Post("/atacar/{id_territorio_origen}/{id_territorio_destino}/{num_dados}", handlers.Atacar)
		r.Post("/ocupar/{territorio_a_ocupar}/{num_ejercitos}", handlers.Ocupar)
		r.Post("/fortificar/{id_territorio_origen}/{id_territorio_destino}/{num_tropas}", handlers.Fortificar)

		// Usuarios
		r.Get("/obtenerNotificaciones", handlers.ObtenerNotificaciones)
		r.Get("/obtenerNumeroNotificaciones", handlers.ObtenerNumeroNotificaciones)
		r.Get("/listarAmigos", handlers.ListarAmigos)
		r.Get("/obtenerPerfil/{nombre}", handlers.ObtenerPerfilUsuario)
		r.Get("/obtenerUsuariosSimilares/{patron}", handlers.ObtenerUsuariosSimilares)
		r.Get("/obtenerSolicitudesPendientes", handlers.ObtenerSolicitudesPendientes)
		r.Get("/consultarTienda", handlers.ConsultarTienda)
		r.Get("/consultarColeccion/{usuario}", handlers.ConsultarColeccion)
		r.Get("/obtenerFotoPerfil/{usuario}", handlers.ObtenerAvatar)
		r.Get("/obtenerDados/{usuario}/{cara}", handlers.ObtenerDados)
		r.Get("/obtenerImagenItem/{id}", handlers.ObtenerImagenItem)
		r.Get("/ranking", handlers.ObtenerRanking)
		r.Get("/eliminarAmigo/{nombre}", handlers.EliminarAmigo)
		r.Post("/aceptarSolicitudAmistad/{nombre}", handlers.AceptarSolicitudAmistad)
		r.Post("/rechazarSolicitudAmistad/{nombre}", handlers.RechazarSolicitudAmistad)
		r.Post("/enviarSolicitudAmistad/{nombre}", handlers.EnviarSolicitudAmistad)
		r.Post("/comprarObjeto/{id_objeto}", handlers.ComprarObjeto)
		r.Post("/modificarBiografia", handlers.ModificarBiografia)
		r.Post("/modificarAspecto/{id_aspecto}", handlers.ModificarAspecto)

	})

	return r
}

// Devuelve un router programado para URLs de cualquiera de los frontends
func routerWeb() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Carpeta del sistema de ficheros que se va a servir, restringida a ella y sus
	// subdirectorios
	filesDir := http.Dir(globales.CARPETA_FRONTEND)
	log.Println("Sirviendo contenido web desde", globales.CARPETA_FRONTEND)

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
