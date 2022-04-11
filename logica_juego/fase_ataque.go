package logica_juego

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

// Ataque permite a un jugador atacar a una región continua, seleccionando el número de dados a utilizar
// TODO completar documentación de la función de ataque
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
		return errors.New("Solo puedes atacar desde un territorio que ocupas")
	}
	if !Conectadas(origen, destino) {
		return errors.New("Solo puedes atacar a un territorio adyacente")
	}
	if atacante == defensor {
		return errors.New("No puedes atacar a un territorio controlado por ti mismo")
	}
	if numDados > 4 || numDados < 1 {
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

	for i >= 0 && j >= 0 {
		if dadosAtacante[i] > dadosDefensor[j] {
			regionDestino.NumTropas--
			tropasPerdidasDefensor++
			if regionDestino.NumTropas == 0 {
				// Región conquistada
				e.HayTerritorioDesocupado = true
				regionDestino.Ocupante = ""

				if e.contarTerritoriosOcupados(defensor) == 0 {
					// Le damos todas las cartas del defensor al atacante
					e.EstadosJugadores[atacante].Cartas = append(e.EstadosJugadores[atacante].Cartas,
						e.EstadosJugadores[defensor].Cartas...)
					e.EstadosJugadores[defensor].Cartas = nil

					// Indicamos que el jugador ha sido derrotado
					e.JugadoresActivos[e.obtenerTurnoJugador(defensor)] = false
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

	// Añadimos la acción correspondiente al ataque
	e.Acciones = append(e.Acciones, NewAccionAtaque(origen, destino, tropasPerdidasAtacante, tropasPerdidasDefensor,
		numDados, atacante, defensor))
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
