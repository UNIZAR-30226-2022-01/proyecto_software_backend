package logica_juego

type Fase int
type TipoTropa int
type NumRegion int

const (
	PUNTOS_PERDER = 10
	PUNTOS_GANAR  = 25
)

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
	NOTIFICACION_PUNTOS
)
