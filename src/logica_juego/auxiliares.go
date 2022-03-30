package logica_juego

func crearSliceJugadores(jugadores []string) (slice []string) {
	for _, jugador := range jugadores {
		slice = append(slice, jugador)
	}

	return slice
}

func crearEstadoMapa() (mapa map[NumRegion]*EstadoRegion) {
	mapa = make(map[NumRegion]*EstadoRegion)
	for i := Eastern_australia; i <= Alberta; i++ {
		mapa[i] = &EstadoRegion{Ocupante: "", NumTropas: 0}
	}

	return mapa
}

func crearMapaEstadosJugadores(jugadores []string) (mapa map[string]*EstadoJugador) {
	mapa = make(map[string]*EstadoJugador, len(jugadores))

	// Tropas iniciales, según el número de jugadores
	numTropas := 0
	if len(jugadores) == 3 {
		numTropas = 35
	} else if len(jugadores) == 4 {
		numTropas = 30
	} else if len(jugadores) == 5 {
		numTropas = 25
	} else if len(jugadores) == 6 {
		numTropas = 20
	}

	for _, j := range jugadores {
		mapa[j] = &EstadoJugador{}

		mapa[j].Cartas = []Carta{}
		mapa[j].UltimoIndiceLeido = -1
		mapa[j].Tropas = numTropas
	}

	return mapa
}
