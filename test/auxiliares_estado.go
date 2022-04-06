package integracion

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

func comprobarAcciones(t *testing.T, cookie *http.Cookie) {
	estado := preguntarEstado(t, cookie)

	if len(estado.Acciones) != (logica_juego.NUM_REGIONES + 2) { // + 2 para que tenga en cuenta cambio de fase
		t.Fatal("Se esperaban", logica_juego.NUM_REGIONES, "acciones en el log, y hay", len(estado.Acciones))
	} else {
		t.Log("Contenidos de acciones:", estado.Acciones)
	}

	estado = preguntarEstado(t, cookie)
	if len(estado.Acciones) != 0 {
		t.Fatal("Se esperaban 0 acciones en el log, y hay", len(estado.Acciones))
	} else {
		t.Log("Contenidos de acciones tras leerlas todas:", estado.Acciones)
	}
}

func preguntarEstado(t *testing.T, cookie *http.Cookie) (estado logica_juego.EstadoPartida) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerEstadoPartida", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de preguntar estado:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al preguntar estado:", resp.StatusCode)
	}

	//body, err := ioutil.ReadAll(resp.Body)
	//bodyString := string(body)
	//t.Log("Respuesta al preguntar estado:", bodyString)

	err = json.NewDecoder(resp.Body).Decode(&estado.Acciones)

	return estado
}

func comprobarConsistenciaEnCurso(t *testing.T, partidaCache vo.Partida) {
	partidaDB := obtenerPartidaDB(t, 1)

	if partidaDB.EnCurso != partidaCache.EnCurso {
		t.Fatal("partidaDB.EnCurso=", partidaDB.EnCurso, "y partidaCache.Encurso=", partidaCache.EnCurso)
	} else {
		if partidaDB.EnCurso {
			t.Log("Ambas partidas en curso")
		} else {
			t.Log("Ambas partidas no en curso")
		}
	}
}

func comprobarConsistenciaAcciones(t *testing.T, partidaCache vo.Partida) {
	time.Sleep(50 * time.Millisecond) // La base de datos no debería tardar mucho más
	partidaDB := obtenerPartidaDB(t, 1)

	if len(partidaDB.Estado.Acciones) != len(partidaCache.Estado.Acciones) {
		t.Fatal("longitud de acciones para partidaDB=", len(partidaDB.Estado.Acciones), ", longitud de acciones para partidaCache=", len(partidaCache.Estado.Acciones))
	} else if !reflect.DeepEqual(partidaDB.Estado.Acciones, partidaCache.Estado.Acciones) {
		t.Fatal("Los estados de la partida en cache y la partida en DB no son consistentes")
	}
}

func obtenerNotificaciones(t *testing.T, cookie *http.Cookie) (notificaciones []interface{}) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerNotificaciones", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en GET de obtener notificaciones:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al obtener notificaciones:", resp.StatusCode)
	} else {
		//body, _ := ioutil.ReadAll(resp.Body)
		//bodyString := string(body)
		//t.Log("Respuesta al obtener notificaciones:", bodyString)

		err = json.NewDecoder(resp.Body).Decode(&notificaciones)
	}

	return notificaciones
}
