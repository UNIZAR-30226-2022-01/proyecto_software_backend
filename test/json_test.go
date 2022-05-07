package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"testing"
)

// Prueba de imprimir acciones en JSON
func TestImpresiónEnJSON(t *testing.T) {
	acciones := make([]interface{}, 9)

	acciones[0] = logica_juego.NewAccionRecibirRegion(1, 4, 8, "usuario1")
	acciones[1] = logica_juego.NewAccionInicioTurno("usuario1", 2, 12, 1)
	acciones[2] = logica_juego.NewAccionCambioFase(2, "usuario1")
	acciones[3] = logica_juego.NewAccionCambioCartas(2, true, []logica_juego.NumRegion{1, 2, 3}, false)
	acciones[4] = logica_juego.NewAccionReforzar("usuario1", 1, 20)
	acciones[5] = logica_juego.NewAccionAtaque(2, 3, 15, 5, []int{2, 2, 3}, []int{3, 4, 6}, "usuario1", "usuario2")
	acciones[6] = logica_juego.NewAccionOcupar(2, 3, 10, 5, "usuario1", "usuario2")
	acciones[7] = logica_juego.NewAccionFortificar(7, 9, 10, 8, "usuario1")
	acciones[8] = logica_juego.NewAccionObtenerCarta(logica_juego.Carta{Tipo: logica_juego.Infanteria, Region: logica_juego.Egypt}, "usuario1")

	for _, a := range acciones {
		serializarAJSONEImprimir(t, a)
	}
}
