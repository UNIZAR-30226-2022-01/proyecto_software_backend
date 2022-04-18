package integracion

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestTienda(t *testing.T) {
	purgarDB()
	cookie := crearUsuario("usuario", t)
	items := consultarTienda(cookie, t)
	if len(items) != 8 {
		t.Fatal("Debería haber 8 objetos disponibles en la tienda")
	}
	t.Log("Se han recuperado los siguientes objetos de la tienda:", items)

	// Le damos puntos al usuario de forma artificial para que pueda comprar
	err := dao.OtorgarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario"}, 100)
	if err != nil {
		t.Fatal("Error al dar puntos:", err)
	}

	items = consultarColeccion(cookie, "usuario", t)
	t.Log(items)
}

func consultarTienda(cookie *http.Cookie, t *testing.T) []vo.ItemTienda {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/consultarTienda", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de consultar tienda:", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		bodyString := string(body)
		t.Fatal("Obtenido código de error no 200 al consultar tienda:", resp.StatusCode, "error:", bodyString)

	} else {
		var items []vo.ItemTienda
		err = json.NewDecoder(resp.Body).Decode(&items)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al consultar tienda:", err)
		}

		return items
	}

	return []vo.ItemTienda{}
}

func consultarColeccion(cookie *http.Cookie, usuario string, t *testing.T) []vo.ItemTienda {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/consultarColeccion/"+usuario, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de consultar coleccion:", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		bodyString := string(body)
		t.Fatal("Obtenido código de error no 200 al consultar coleccion:", resp.StatusCode, "error:", bodyString)

	} else {
		var items []vo.ItemTienda
		err = json.NewDecoder(resp.Body).Decode(&items)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al consultar coleccion:", err)
		}

		return items
	}

	return []vo.ItemTienda{}
}
