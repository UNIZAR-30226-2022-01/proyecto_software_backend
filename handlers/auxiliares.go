package handlers

import (
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"sort"
)

// Funciones auxiliares

func ordenarPorNumeroJugadores(writer http.ResponseWriter, partidasPrivadasSinAmigos []vo.Partida) {
	sort.SliceStable(partidasPrivadasSinAmigos, func(i, j int) bool {
		// Orden: > a <
		jugadoresI, _, err1 := dao.ConsultarJugadoresPartida(globales.Db, &partidasPrivadasSinAmigos[i])
		jugadoresJ, _, err2 := dao.ConsultarJugadoresPartida(globales.Db, &partidasPrivadasSinAmigos[j])

		if err1 != nil || err2 != nil {
			devolverErrorSQL(writer)
		}

		return len(jugadoresI) > len(jugadoresJ)
	})
}

func dividirPartidasPorAmigos(partidasPrivadas []vo.Partida, amigos []vo.Usuario) ([]vo.Partida, []vo.Partida) {
	var partidasPrivadasConAmigos []vo.Partida
	var partidasPrivadasSinAmigos []vo.Partida
	for _, partida := range partidasPrivadas {
		// Se ha llegado al punto en el slice a partir del cual no hay amigos
		jugadores, _, _ := dao.ConsultarJugadoresPartida(globales.Db, &partida)

		amigos := obtenerAmigos(amigos, jugadores)

		if len(amigos) == 0 {
			partidasPrivadasSinAmigos = append(partidasPrivadasSinAmigos, partida)
		} else {
			partidasPrivadasConAmigos = append(partidasPrivadasConAmigos, partida)
		}
	}
	return partidasPrivadasConAmigos, partidasPrivadasSinAmigos
}

func ordenarPorNumeroAmigos(partidasPrivadas []vo.Partida, amigos []vo.Usuario) {
	sort.SliceStable(partidasPrivadas, func(i, j int) bool {
		jugadores, _, _ := dao.ConsultarJugadoresPartida(globales.Db, &partidasPrivadas[i])
		listaAmigosI := obtenerAmigos(amigos, jugadores)

		jugadores, _, _ = dao.ConsultarJugadoresPartida(globales.Db, &partidasPrivadas[j])
		listaAmigosJ := obtenerAmigos(amigos, jugadores)

		// Orden: > a <
		return len(listaAmigosI) > len(listaAmigosJ)
	})
}

func dividirPartidasPrivadasYPublicas(partidas []vo.Partida) ([]vo.Partida, []vo.Partida) {
	// Extrae las partidas privadas del slice y deja las partidas públicas
	var partidasPrivadas []vo.Partida
	var partidasPublicas []vo.Partida
	for _, partida := range partidas {
		if !partida.EsPublica {
			partidasPrivadas = append(partidasPrivadas, partida)
		} else {
			partidasPublicas = append(partidasPublicas, partida)
		}
	}
	return partidasPrivadas, partidasPublicas
}

// transformarAElementoListaPartidas convierte una partida en un elemento de lista de partidas,
// dada una lista de amigos de un usuario. Se asume que la partida existe en la DB.
// No puede localizarse en el módulo VO porque causaría una dependencia cíclica con DAO
func transformarAElementoListaPartidas(p *vo.Partida, amigos []vo.Usuario) vo.ElementoListaPartidas {
	jugadores, _, _ := dao.ConsultarJugadoresPartida(globales.Db, p)
	listaAmigos := obtenerAmigos(amigos, jugadores)

	return vo.ElementoListaPartidas{
		IdPartida:          p.IdPartida,
		EsPublica:          p.EsPublica,
		NumeroJugadores:    len(jugadores),
		MaxNumeroJugadores: p.MaxNumeroJugadores,
		AmigosPresentes:    listaAmigos,
		NumAmigosPresentes: len(listaAmigos),
	}
}

// obtenerAmigos obtiene una lista de nombres amigos presentes en
// una partida, dada una lista previa
func obtenerAmigos(amigos []vo.Usuario, jugadores []vo.Usuario) (listaFiltrada []string) {
	for _, amigo := range amigos {
		// Como máximo hay 6 jugadores en la partida, así que
		// la complejidad la dicta el número de amigos del usuario
		for _, jugador := range jugadores {
			if amigo.NombreUsuario == jugador.NombreUsuario {
				listaFiltrada = append(listaFiltrada, amigo.NombreUsuario)
			}
		}
	}

	return listaFiltrada
}

// hashPassword crea un hash de clave utilizando bcrypt
// https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) // Coste fijo por defecto para evitar tiempos de cálculo excesivos
	return string(bytes), err
}

// Devuelve una respuesta con status 500 junto al mensaje de error y la función
// en la que se ha dado.
func devolverError(writer http.ResponseWriter, err error) {
	log.Println("Error:", err)
	writer.WriteHeader(http.StatusInternalServerError)
	_, err = writer.Write([]byte(err.Error()))
	if err != nil {
		log.Println("Error al escribir respuesta en:", err)
	}
}

// Devuelve una respuesta con status 500 junto al mensaje de error y la función
// en la que se ha dado.
func devolverErrorSQL(writer http.ResponseWriter) {
	err := errors.New("Se ha producido un error en la base de datos.")

	log.Println("Error en:", err)
	writer.WriteHeader(http.StatusInternalServerError)
	_, err = writer.Write([]byte(err.Error()))
	if err != nil {
		log.Println("Error al escribir respuesta en:", err)
	}
}

// Devuelve una respuesta con status 200.
func escribirHeaderExito(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
}

func transformaAElementoListaUsuarios(usuario vo.Usuario) vo.ElementoListaUsuarios {
	return vo.ElementoListaUsuarios{
		NombreUsuario:   usuario.NombreUsuario,
		Email:           usuario.Email,
		Biografia:       usuario.Biografia,
		PartidasGanadas: usuario.PartidasGanadas,
		PartidasTotales: usuario.PartidasTotales,
		Puntos:          usuario.Puntos,
		ID_dado:         usuario.ID_dado,
		ID_avatar:       usuario.ID_avatar,
		EsAmigo:         false,
	}
}
