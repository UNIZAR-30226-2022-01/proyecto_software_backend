package logica_juego

const (
	Eastern_australia NumRegion = iota
	Indonesia
	New_guinea
	Alaska
	Ontario
	Northwest_territory
	Venezuela
	Madagascar
	North_africa
	Greenland
	Iceland
	Great_britain
	Scandinavia
	Japan
	Yakursk
	Kamchatka
	Siberia
	Ural
	Afghanistan
	Middle_east
	India
	Siam
	China
	Mongolia
	Irkutsk
	Ukraine
	Southern_europe
	Western_europe
	Northern_europe
	Egypt
	East_africa
	Congo
	South_africa
	Brazil
	Argentina
	Eastern_united_states
	Western_united_states
	Quebec
	Central_america
	Peru
	Western_australia
	Alberta
)

func (nr NumRegion) String() string {
	// Lo sentimos
	return []string{"eastern_australia", "indonesia",
		"new_guinea", "alaska", "ontario", "northwest_territory",
		"venezuela", "madagascar", "north_africa", "greenland",
		"iceland", "great_britain", "scandinavia", "japan", "yakursk",
		"kamchatka", "siberia", "ural", "afghanistan", "middle_east",
		"india", "siam", "china", "mongolia", "irkutsk", "ukraine",
		"southern_europe", "western_europe", "northern_europe", "egypt",
		"east_africa", "congo", "south_africa", "brazil", "argentina",
		"eastern_united_states", "western_united_states", "quebec",
		"central_america", "peru", "western_australia", "alberta"}[nr]
}

// Continente define el tipo utilizado para representar cada continente. Se almacenará el valor de dicho continente,
// que corresponde con el número de tropas de bonificación recibidas al ocuparlo, además de la lista de regiones que
// lo compone
type Continente struct {
	Valor    int
	Regiones []NumRegion
}

// Continentes que componen el mapa del juego
var Continentes map[string]Continente

// InicializarContinentes inicializa el mapa Continentes con cada uno de los continentes del mapa de juego
func InicializarContinentes() {
	Continentes = make(map[string]Continente)

	Continentes["América del Norte"] = Continente{
		Valor: 5,
		Regiones: []NumRegion{Alaska, Northwest_territory, Greenland, Alberta, Ontario,
			Quebec, Western_united_states, Eastern_united_states, Central_america}}

	Continentes["Europa"] = Continente{
		Valor: 5,
		Regiones: []NumRegion{Iceland, Scandinavia, Great_britain, Northern_europe,
			Ukraine, Western_europe, Southern_europe}}

	Continentes["Asia"] = Continente{
		Valor: 7,
		Regiones: []NumRegion{Yakursk, Ural, Siberia, Irkutsk, Kamchatka, Afghanistan,
			China, Mongolia, Japan, Middle_east, India, Siam}}

	Continentes["América del Sur"] = Continente{
		Valor:    2,
		Regiones: []NumRegion{Venezuela, Brazil, Peru, Argentina}}

	Continentes["África"] = Continente{
		Valor:    3,
		Regiones: []NumRegion{North_africa, Egypt, Congo, East_africa, South_africa, Madagascar}}

	Continentes["Oceanía"] = Continente{
		Valor:    2,
		Regiones: []NumRegion{Indonesia, New_guinea, Western_australia, Eastern_australia}}
}

// ContarTerritoriosOcupados cuenta el número de territorios que ocupa un jugador determinado
func (e *EstadoPartida) ContarTerritoriosOcupados(jugador string) int {
	n := 0
	for i := Eastern_australia; i <= Alberta; i++ {
		if e.EstadoMapa[i].Ocupante == jugador {
			n++
		}
	}
	return n
}
