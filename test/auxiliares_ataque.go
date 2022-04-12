package integracion

import (
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/logica_juego"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func atacar(origen, destino logica_juego.NumRegion, numDados int, cookie *http.Cookie, t *testing.T) error {
	client := &http.Client{}
	idOrigen := int(origen)
	idDestino := int(destino)
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/atacar/"+
		strconv.Itoa(idOrigen)+"/"+strconv.Itoa(idDestino)+"/"+strconv.Itoa(numDados), nil)
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
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		return errors.New(bodyString)
	}

	return nil
}

func ocupar(territorioAOcupar logica_juego.NumRegion, numTropas int, cookie *http.Cookie, t *testing.T) error {
	client := &http.Client{}
	idRegion := int(territorioAOcupar)
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/ocupar/"+
		strconv.Itoa(idRegion)+"/"+strconv.Itoa(numTropas), nil)
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
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		return errors.New(bodyString)
	}

	return nil
}
