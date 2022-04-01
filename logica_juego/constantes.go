package logica_juego

type Fase int
type TipoTropa int
type NumRegion int

const (
	Inicio Fase = iota // Repartir regiones
	Refuerzo
	Ataque
	Fortificar
)

const (
	Infanteria TipoTropa = iota
	Caballeria
	Artilleria
)

const (
	NUM_REGIONES = 42
)

const (
	NOTIFICACION_AMISTAD = iota
	NOTIFICACION_TURNO
)
