package integracion

import (
	"encoding/json"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/middleware"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/vo"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func crearUsuario(nombre string, t *testing.T) (cookie *http.Cookie) {
	cookie = nil

	campos := url.Values{
		"nombre":   {nombre},
		"email":    {nombre + "@" + nombre + ".com"},
		"password": {nombre},
	}
	resp, err := http.PostForm("http://localhost:"+os.Getenv(globales.PUERTO_API)+"/registro", campos)
	if err != nil {
		t.Fatal("No se ha podido realizar request POST:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al registrar un usuario:", resp.StatusCode)
	}

	// De ObtenerUsuarioCookie
	for _, c := range resp.Cookies() {
		if c.Name == middleware.NOMBRE_COOKIE_USUARIO { // Es una cookie de usuario
			// Obtener el usuario del valor de la cookie
			nombreCookie := c.Value[:strings.IndexRune(c.Value, middleware.SEPARADOR_VALOR_COOKIE_USUARIO)]
			if nombre != nombreCookie {
				t.Fatal("Obtenido nombre de cookie diferente del esperado:", nombreCookie, "esperaba:", nombre)
			}
			cookie = c
			break
		}
	}

	if cookie == nil {
		t.Fatal("No se ha obtenido una cookie en la respuesta de crear usuario para", nombre)
	}

	return cookie
}

func solicitarAmistad(cookie *http.Cookie, t *testing.T, nombre string) {
	t.Log("Solicitando amistad de userPrincipal a", nombre)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/enviarSolicitudAmistad/"+nombre, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de solicitar amistad:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al solicitar amistad:", resp.StatusCode)
	}
}

func solicitarAmistadConError(cookie *http.Cookie, t *testing.T, nombre string) {
	t.Log("Solicitando amistad de a", nombre)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/enviarSolicitudAmistad/"+nombre, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de solicitar amistad:", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatal("Se esperaba obtener un error:", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	t.Log("OK se ha obtenido un error al realizar la solicitud de amistad:", string(body))
}

func aceptarSolicitudDeAmistad(cookie *http.Cookie, t *testing.T, nombre string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/aceptarSolicitudAmistad/"+nombre, nil) // MAPS :D
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de solicitar amistad:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al aceptar amistad:", resp.StatusCode)
	}
}

func rechazarSolicitudDeAmistad(cookie *http.Cookie, t *testing.T, nombre string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/rechazarSolicitudAmistad/"+nombre, nil) // MAPS :D
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de rechazar amistad:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al rechazar amistad:", resp.StatusCode)
	}
}

func listarAmigos(cookie *http.Cookie, t *testing.T) []string {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/listarAmigos", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de listar amigos:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al listar amigos:", resp.StatusCode)
	} else {
		var amigos []string
		err = json.NewDecoder(resp.Body).Decode(&amigos)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al listar amigos:", err)
		}

		t.Log("Respuesta de listarAmigos:", amigos)
		return amigos
	}

	return nil
}

func consultarSolicitudesPendientes(cookie *http.Cookie, t *testing.T) []string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerSolicitudesPendientes", nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("Error en GET de consultar amigos pendientes:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al consultar pendientes:", resp.StatusCode)
	} else {
		var pendientes []string
		err = json.NewDecoder(resp.Body).Decode(&pendientes)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al listar pendientes:", err)
		}

		t.Log("Respuesta de consultarPendientes:", pendientes)
		return pendientes
	}

	return []string{}
}

func obtenerPerfilUsuario(cookie *http.Cookie, nombre string, t *testing.T) vo.ElementoListaUsuarios {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerPerfil/"+nombre, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de consultar perfil:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al obtener perfil:", resp.StatusCode)
	} else {
		var usuario vo.ElementoListaUsuarios
		err = json.NewDecoder(resp.Body).Decode(&usuario)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al obtener perfil:", err)
		}

		return usuario
	}

	return vo.ElementoListaUsuarios{}
}

func buscarUsuariosSimilares(cookie *http.Cookie, patron string, t *testing.T) []vo.ElementoListaUsuariosSimilares {
	cliente := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/obtenerUsuariosSimilares/"+patron, nil)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	resp, err := cliente.Do(req)
	if err != nil {
		t.Fatal("Error en GET de buscar usuarios:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al buscar usuarios:", resp.StatusCode)
	} else {
		var usuarios []vo.ElementoListaUsuariosSimilares

		err = json.NewDecoder(resp.Body).Decode(&usuarios)
		if err != nil {
			t.Fatal("Error al leer JSON de respuesta al buscar usuarios:", err)
		}

		return usuarios
	}

	return nil
}

func modificarBiografia(cookie *http.Cookie, biografia string, t *testing.T) {
	client := &http.Client{}

	var campos url.Values
	campos = url.Values{
		"biografia": {biografia},
	}

	req, err := http.NewRequest("POST", "http://localhost:"+os.Getenv(globales.PUERTO_API)+"/api/modificarBiografia", strings.NewReader(campos.Encode()))
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Para indicar que el formulario "va en la url", porque campos.Encode() hace eso

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal("Error en POST de modificar biografia:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al modificar biografia:", resp.StatusCode)
	}
}
