package logica_juego

import (
	"errors"
	"log"
	"math/rand"
	"sort"
	"time"
)

// Ataque permite a un usuario atacar a un territorio adyacente a alguna de sus regiones. Para ello, debe de tener por
// lo menos dos ejércitos en el territorio desde el que ataca. Al atacar deberá elegir el número de dados a lanzar, entre
// 1 y 3. Cabe destacar que será necesario tener al menos un ejército más que el número de dados a lanzar, por ejemplo, si
// quiero lanzar 3 dados, el territorio tendrá que tener 4 ejércitos por lo menos.
// Por otro lado el defensor tirará 2 dados si tiene 2 ejércitos o más, o 1 en el caso contrario.
//
// Para calcular el resultado del ataque, se compararán los dados con mayor valor de ambos jugadores. Si el atacante consigue
// un resultado mayor, el territorio defensor perderá una tropa. Por otro lado, si empatan o gana el defensor el territorio
// atacante perderá un ejército. En caso de que ambos jugadores hayan lanzado más de un dado, se repetirá el mismo proceso
// comparando el valor del segundo dado más alto de cada uno
//
// No se puede atacar en los siguientes casos: no es el turno deñ jugador, no es la fase de ataque, el jugador tiene más de
// 4 cartas, hay algún territorio sin ocupar, el territorio atacado no es adyacente, el territorio atacado no es de un rival,
// el número de dados no está entre 1 y 3 o el número de ejércitos no supera el número de dados.
func (e *EstadoPartida) Ataque(origen, destino NumRegion, numDados int, jugador string) error {
	regionOrigen := e.EstadoMapa[origen]
	regionDestino := e.EstadoMapa[destino]
	atacante := regionOrigen.Ocupante
	defensor := regionDestino.Ocupante

	if e.ObtenerJugadorTurno() != jugador {
		return errors.New("Solo puedes atacar durante tu turno")
	}
	if e.Fase != Ataque {
		return errors.New("Solo puedes atacar durante la fase de ataque")
	}
	if len(e.EstadosJugadores[jugador].Cartas) >= 5 {
		return errors.New("Estás obligado a cambiar cartas si tienes 5 o más")
	}
	if e.HayTerritorioDesocupado {
		return errors.New("No puedes atacar si hay algún territorio sin ocupar")
	}
	if jugador != atacante {
		log.Println("DEVUELVO ERROR TERRITORIO PROPIO")
		return errors.New("Solo puedes atacar desde un territorio que ocupas")
	}
	if !Conectadas(origen, destino) {
		return errors.New("Solo puedes atacar a un territorio adyacente")
	}
	if atacante == defensor {
		return errors.New("No puedes atacar a un territorio controlado por ti mismo")
	}
	if numDados > 3 || numDados < 1 {
		return errors.New("Solo puedes lanzar 1, 2 o 3 dados")
	}
	if numDados >= regionOrigen.NumTropas {
		return errors.New("Necesitas al menos un ejército más que el número de dados a lanzar")
	}

	numDadosRival := 1
	if regionDestino.NumTropas > 1 {
		// El defensor siempre lanza 2 dados si tiene más de un ejército
		numDadosRival = 2
	}

	// Lanzamos los dados y los ordenamos de menor a mayor
	dadosAtacante := lanzarNDados(numDados)
	dadosDefensor := lanzarNDados(numDadosRival)
	sort.Sort(sort.IntSlice(dadosAtacante))
	sort.Sort(sort.IntSlice(dadosDefensor))

	i := numDados - 1
	j := numDadosRival - 1
	tropasPerdidasAtacante := 0
	tropasPerdidasDefensor := 0

	emitirAccionJugadorEliminado := false
	numCartasDefensor := 0
	for i >= 0 && j >= 0 {
		if dadosAtacante[i] > dadosDefensor[j] {
			regionDestino.NumTropas--
			tropasPerdidasDefensor++
			if regionDestino.NumTropas == 0 {
				// Región conquistada
				e.HayTerritorioDesocupado = true
				regionDestino.Ocupante = ""

				if e.ContarTerritoriosOcupados(defensor) == 0 {
					// Le damos todas las cartas del defensor al atacante
					numCartasDefensor = len(e.EstadosJugadores[defensor].Cartas)

					e.EstadosJugadores[atacante].Cartas = append(e.EstadosJugadores[atacante].Cartas,
						e.EstadosJugadores[defensor].Cartas...)
					e.EstadosJugadores[defensor].Cartas = nil

					// Indicamos que el jugador ha sido derrotado
					e.JugadoresActivos[e.obtenerTurnoJugador(defensor)] = false

					emitirAccionJugadorEliminado = true // Emitir acción de ataque después de eliminarlo
				}
				break
			}
		} else {
			regionOrigen.NumTropas--
			tropasPerdidasAtacante++
			if regionOrigen.NumTropas <= 1 {
				// El atacante no sigue atacando si tiene 1 ejército
				break
			}
		}
		i--
		j--
	}

	// Actualizamos el estado del último ataque
	e.DadosUltimoAtaque = numDados
	e.TropasPerdidasUltimoAtaque = tropasPerdidasAtacante
	e.RegionUltimoAtaque = origen
	e.UltimoDefensor = defensor

	// Añadimos la acción correspondiente al ataque
	e.Acciones = append(e.Acciones, NewAccionAtaque(origen, destino, tropasPerdidasAtacante, tropasPerdidasDefensor,
		numDados, atacante, defensor))

	// Se emite después del ataque
	if emitirAccionJugadorEliminado {
		e.Acciones = append(e.Acciones, NewAccionJugadorEliminado(defensor, atacante, numCartasDefensor))
	}

	return nil
}

// Ocupar permite a un usuario ocupar un territorio sin tropas, especificando el territorio a ocupar y el número de
// tropas que quiere mover a él. Dichas tropas se moverán desde el territorio con el que conquistó la región a ser ocupada.
//
// Para ocupar se deben cumplir las siguientes condiciones: hay alguna región sin tropas, dicha región es adyacente a
// la región desde la que se inició el último ataque, la ocupación se realiza durante el turno del jugador y en la fase
// de ataque, el número de tropas asignadas por la ocupación no deja al territorio origen sin tropas, el número de tropas
// asignadas es mayor al número de dados usados en el último ataque menos el número de ejércitos que perdió el atacante
// en dicho ataque.
//
// Cabe destacar que siempre que un territorio quede sin tropas tras un ataque, el juego no permitirá continuar atacando
// ni cambiar de fase o turno hasta que dicho territorio sea ocupado, de manera que solo podrá haber un territorio
// sin ocupar a la vez.
func (e *EstadoPartida) Ocupar(territorio NumRegion, numEjercitos int, jugador string) error {
	// Comprobación de errores
	if e.ObtenerJugadorTurno() != jugador {
		return errors.New("No puedes ocupar un territorio fuera de tu turno")
	}
	if e.Fase != Ataque {
		return errors.New("No puedes ocupar fuera de la fase de ataque")
	}
	if len(e.EstadosJugadores[jugador].Cartas) > 4 {
		return errors.New("No puedes ocupar un territorio si tienes más de 4 cartas")
	}
	if !e.HayTerritorioDesocupado {
		return errors.New("No se puede ocupar si no hay territorios desocupados")
	}
	if e.EstadoMapa[territorio].NumTropas > 0 {
		return errors.New("No se puede ocupar un territorio con tropas")
	}
	if !Conectadas(e.RegionUltimoAtaque, territorio) {
		return errors.New("No puedes ocupar un territorio desde una región no adyacente")
	}
	if numEjercitos < e.DadosUltimoAtaque-e.TropasPerdidasUltimoAtaque {
		return errors.New("Debes ocupar el territorio con al menos el número de dados usados en el último ataque," +
			"menos el número de tropas pérdidas en dicho ataque")
	}
	if numEjercitos >= e.EstadoMapa[e.RegionUltimoAtaque].NumTropas {
		return errors.New("No puedes dejar al territorio desde el que ocupas sin tropas")
	}

	// Ocupamos el territorio
	e.HayTerritorioDesocupado = false
	e.EstadoMapa[territorio].Ocupante = jugador
	e.EstadoMapa[territorio].NumTropas = numEjercitos
	e.EstadoMapa[e.RegionUltimoAtaque].NumTropas -= numEjercitos

	// Comprobamos si ha ganado la partida
	if e.ContarTerritoriosOcupados(jugador) == NUM_REGIONES {
		// Se introduce una acción de fin de partida
		e.Acciones = append(e.Acciones, NewAccionPartidaFinalizada(jugador))

		// Y se marca como terminada, para borrarla tras ser consultada por todos los jugadores
		e.Terminada = true
	}

	// Añadimos la acción de ocupación
	e.Acciones = append(e.Acciones, NewAccionOcupar(e.RegionUltimoAtaque, territorio,
		e.EstadoMapa[e.RegionUltimoAtaque].NumTropas, e.EstadoMapa[territorio].NumTropas,
		jugador, e.UltimoDefensor))

	return nil
}

// lanzarNDados simula el lanzamiento de un número determinado "n" de dados
// Los resultados de dichos lanzamientos son devueltos como un slice de enteros
func lanzarNDados(n int) (dados []int) {
	onlyOnce.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	for i := 0; i < n; i++ {
		dados = append(dados, rand.Intn(6)+1)
	}

	return dados
}
