package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func reforzarTerritorioConFallo(t *testing.T, cookie *http.Cookie, numRegion int, numTropas int) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/reforzarTerritorio/"+strconv.Itoa(numRegion)+"/"+strconv.Itoa(numTropas), nil)
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
		t.Fatal("Obtenido código de error OK al forzar error en reforzar territorio:", resp.StatusCode)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		t.Log("Recibido error correctamente: " + bodyString)
	}

}

func reforzarTerritorio(t *testing.T, cookie *http.Cookie, numRegion int, numTropas int) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/reforzarTerritorio/"+strconv.Itoa(numRegion)+"/"+strconv.Itoa(numTropas), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace e

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en POST de reforzar territorio:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al reforzar territorio:", resp.StatusCode)
	}

}
