package integracion

import (
	"encoding/json"
	"errors"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestTienda(t *testing.T) {
	purgarDB()
	var err error
	cookie := crearUsuario("usuario", t)
	items := consultarTienda(cookie, t)
	if len(items) != 8 {
		t.Fatal("Debería haber 8 objetos disponibles en la tienda")
	}
	t.Log("Se han recuperado los siguientes objetos de la tienda:", items)

	// Le damos puntos al usuario de forma artificial para que pueda comprar
	dao.OtorgarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario"}, 100)

	// Intentamos comprar un objeto inexistente, se espera error
	t.Log("Intentamos comprar un objeto inexistente, se espera error")
	err = comprarObjeto(cookie, 500, t)
	if err == nil {
		t.Fatal("Se esperaba error al comprar un objeto inexistente")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Intentamos comprar un objeto con puntos insuficientes, se espera error
	t.Log("Intentamos comprar un objeto inexistente, se espera error")
	dao.RetirarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario"}, 100)
	err = comprarObjeto(cookie, 5, t)
	if err == nil {
		t.Fatal("Se esperaba error al comprar un objeto inexistente")
	}
	t.Log("OK, se ha obtenido el error:", err)
	dao.OtorgarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario"}, 100)

	// Intentamos comprar un objeto inicial, se espera error
	t.Log("Intentamos comprar un objeto inicial, se espera error")
	err = comprarObjeto(cookie, 0, t)
	if err == nil {
		t.Fatal("Se esperaba error al comprar un objeto inicial")
	}
	t.Log("OK, se ha obtenido el error:", err)

	// Comprar objeto correcto
	t.Log("Intentamos comprar un objeto existente y con suficiente dinero")
	err = comprarObjeto(cookie, 5, t)
	if err != nil {
		t.Fatal("Error al comprar objeto:", err)
	}
	t.Log("Se ha comprado el objeto correctamente")

	// Intentamos comprar un objeto que ya tenemos, se espera error
	t.Log("Intentamos comprar un objeto que ya tenemos, se espera error")
	err = comprarObjeto(cookie, 5, t)
	if err == nil {
		t.Fatal("Se esperaba error al comprar un objeto que ya tenemos")
	}
	t.Log("OK, se ha obtenido el error:", err)

	items = consultarColeccion(cookie, "usuario", t)
	if len(items) != 1 {
		t.Fatal("No se ha consultado correctamente la colección de objetos del jugador")
	}
	t.Log("Se han recuperado los siguientes objetos del jugador:", items)
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

func comprarObjeto(cookie *http.Cookie, idObjeto int, t *testing.T) error {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/comprarObjeto/"+strconv.Itoa(idObjeto), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de comprar objeto:", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		bodyString := string(body)
		return errors.New(bodyString)
	}

	return nil
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
