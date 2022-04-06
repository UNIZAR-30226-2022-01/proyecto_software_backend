package integracion

import (
	"encoding/json"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/servidor"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func purgarDB() {
	_, err := globales.Db.Exec(`DELETE FROM "backend"."Partida"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."EsAmigo"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."Participa"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."TieneItems"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`DELETE FROM "backend"."Usuario"`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = globales.Db.Exec(`ALTER SEQUENCE backend."Partida_id_seq" RESTART`)
	if err != nil {
		log.Fatal(err)
	}
}

// Función que se ejecuta automáticamente antes de los test
func init() {
	// Inyecta las variables de entorno
	os.Setenv("DIRECCION_DB", "postgres")
	os.Setenv("DIRECCION_DB_TESTS", "localhost")
	os.Setenv("PUERTO_API", "8090")
	os.Setenv("PUERTO_WEB", "8080")
	os.Setenv("USUARIO_DB", "postgres")
	os.Setenv("PASSWORD_DB", "postgres")

	go servidor.IniciarServidor(true)
	time.Sleep(5 * time.Second)
}

// Funciones auxiliares generales

func serializarAJSONEImprimir(t *testing.T, obj interface{}) {
	bytes, err := json.MarshalIndent(obj, "", "\t")

	if err != nil {
		t.Fatal("Error al serializar a JSON", obj, ":", err)
	} else {
		t.Log("JSON de", obj, ":", string(bytes))
	}
}

func saltarFase(cookie *http.Cookie, t *testing.T) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/pasarDeFase", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en GET de obtener partidas:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Log("No se ha podido saltar la fase")
		body, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(body)
		return errors.New(bodyString)
	}

	return nil
}
