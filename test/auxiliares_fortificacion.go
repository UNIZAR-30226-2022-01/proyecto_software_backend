package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func fortificarTerritorioConError(t *testing.T, cookie *http.Cookie, numTropas int, territorio1 int, territorio2 int) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/fortificarTerritorio/"+strconv.Itoa(territorio1)+"/"+strconv.Itoa(territorio2)+"/"+strconv.Itoa(numTropas), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace e

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en POST de reforzar territorio:", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatal("Obtenido código de error OK al forzar error en fortificar territorio:", resp.StatusCode)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		t.Log("Recibido error correctamente: " + bodyString)
	}
}

func fortificarTerritorio(t *testing.T, cookie *http.Cookie, numTropas int, territorio1 int, territorio2 int) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/fortificarTerritorio/"+strconv.Itoa(territorio1)+"/"+strconv.Itoa(territorio2)+"/"+strconv.Itoa(numTropas), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace e

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en POST de fortificar territorio:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al fortificar territorio:", resp.StatusCode)
	}
}

// Fuerza el paso a la fase de fortificar, sin comprobaciones
func pasarAFaseFortificar(partidaCache vo.Partida) {
	partidaCache.Estado.Fase = logica_juego.Fortificar
	globales.CachePartidas.AlmacenarPartida(partidaCache)
}
