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

type Continente struct {
	Valor    int
	Regiones []NumRegion
}

var Continentes map[string]Continente

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
