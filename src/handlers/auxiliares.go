package handlers

import (
	"backend/dao"
	"backend/globales"
	"backend/vo"
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
