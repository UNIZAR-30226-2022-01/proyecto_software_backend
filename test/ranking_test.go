package integracion

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"net/http"
	"os"
	"testing"
)

func TestRanking(t *testing.T) {
	purgarDB()
	cookie := crearUsuario("usuario", t)

	// Añadimos usuarios con distintos números de victorias
	_, err := globales.Db.Exec(`INSERT INTO backend."Usuario" (email, "nombreUsuario", "passwordHash", biografia, "cookieSesion",
                               "partidasGanadas", "partidasTotales", puntos, "ID_avatar", "ID_dado") VALUES 
                               ('email1', 'juan', 'password', 'biografia', 'cookie', 15, 20, 100, 1, 9),
                               ('email2', 'pedro', 'password', 'biografia', 'cookie', 10, 12, 100, 1, 9),
                               ('email3', 'fran', 'password', 'biografia', 'cookie', 30, 35, 100, 1, 9),
                               ('email4', 'luis', 'password', 'biografia', 'cookie', 5, 10, 100, 1, 9),
                               ('email5', 'susana', 'password', 'biografia', 'cookie', 6, 15, 100, 1, 9),
                               ('email6', 'daniel', 'password', 'biografia', 'cookie', 15, 20, 100, 1, 9),
                               ('email7', 'usuario1', 'password', 'biografia', 'cookie', 17, 20, 100, 1, 9),
                               ('email8', 'usuario2', 'password', 'biografia', 'cookie', 29, 50, 100, 1, 9);`)
	if err != nil {
		t.Fatal("Error al crear los usuarios:", err)
	}

	// Comprobamos que el ranking esté en orden
	ranking := obtenerRanking(cookie, t)
	if len(ranking) != 9 {
		t.Fatal("No todos los jugadores aparecen en el ranking")
	}
	for i, u := range ranking {
		if u.PartidasGanadas < ranking[i+1].PartidasGanadas {
			t.Fatal("El ranking no está ordenado correctamente")
		}
		if i == len(ranking)-2 {
			break
		}
	}

	t.Log("El ranking está ordenado correctamente")
	for i, u := range ranking {
		t.Log("Jugador", i, ", nombre:", u.NombreUsuario, "partidas ganadas:", u.PartidasGanadas)
	}
}

func obtenerRanking(cookie *http.Cookie, t *testing.T) []vo.ElementoRankingUsuarios {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/ranking", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de obtener ranking:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al obtener ranking:", resp.StatusCode)
	} else {
		var ranking []vo.ElementoRankingUsuarios

		err = json.NewDecoder(resp.Body).Decode(&ranking)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al obtener ranking:", err)
		}

		return ranking
	}

	return nil
}
