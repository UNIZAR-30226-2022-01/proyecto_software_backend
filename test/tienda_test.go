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
	dao.OtorgarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario"}, 100, true)

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
	dao.OtorgarPuntos(globales.Db, &vo.Usuario{NombreUsuario: "usuario"}, 100, true)

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

	t.Log("Intentamos comprar un objeto existente y con suficiente dinero")
	err = comprarObjeto(cookie, 7, t)
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
	if len(items) != 2 {
		t.Fatal("No se ha consultado correctamente la colección de objetos del jugador")
	}
	t.Log("Se han recuperado los siguientes objetos del jugador:", items)

	// Intentamos equipar el aspecto de dados comprado
	t.Log("Intentamos equipar el aspecto de dados comprado")
	err = modificarAspecto(cookie, 7, t)
	if err != nil {
		t.Fatal("Error al modificar el aspecto:", err)
	}
	usuario := obtenerPerfilUsuario(cookie, "usuario", t)
	if usuario.ID_dado != 7 {
		t.Fatal("No se ha equipado el dado correctamente")
	}
	t.Log("Se ha equipado el dado correctamente")

	// Intentamos equipar el aspecto de fichas comprado
	t.Log("Intentamos equipar el aspecto de fichas comprado")
	err = modificarAspecto(cookie, 5, t)
	if err != nil {
		t.Fatal("Error al modificar el aspecto:", err)
	}
	usuario = obtenerPerfilUsuario(cookie, "usuario", t)
	if usuario.ID_ficha != 5 {
		t.Fatal("No se ha equipado el aspecto correctamente")
	}
	t.Log("Se ha equipado el aspecto de ficha correctamente")

	// Intentamos equipar un aspecto por defecto
	t.Log("Intentamos equipar un aspecto por defecto")
	err = modificarAspecto(cookie, 0, t)
	if err != nil {
		t.Fatal("Error al modificar el aspecto:", err)
	}
	usuario = obtenerPerfilUsuario(cookie, "usuario", t)
	if usuario.ID_ficha != 0 {
		t.Fatal("No se ha equipado el aspecto correctamente")
	}
	t.Log("Se ha equipado el aspecto correctamente")

	// Intentamos equipar un aspecto inexistente, se espera error
	t.Log("Intentamos equipar un aspecto inexistente, se espera error")
	err = modificarAspecto(cookie, 200, t)
	if err == nil {
		t.Fatal("Se esperaba error al intentar equipar un aspecto inexistente")
	}
	t.Log("OK, error obtenido:", err)

	// Intentamos equipar un aspecto que no tiene el jugador, se espera error
	t.Log("Intentamos equipar un aspecto que no tiene el jugador, se espera error")
	err = modificarAspecto(cookie, 8, t)
	if err == nil {
		t.Fatal("Se esperaba error al intentar equipar un aspecto que no tiene el jugador")
	}
	t.Log("OK, error obtenido:", err)
}

func modificarAspecto(cookie *http.Cookie, idAspecto int, t *testing.T) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/modificarAspecto/"+strconv.Itoa(idAspecto), nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de modificar aspecto:", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		bodyString := string(body)
		return errors.New(bodyString)
	}

	return nil
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
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/comprarObjeto/"+strconv.Itoa(idObjeto), nil)
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
