package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func login(t *testing.T, usuario, pass string) {
	campos := url.Values{
		"nombre":   {usuario},
		"password": {pass},
	}

	resp, err := http.PostForm("http://localhost:"+os.Getenv(globales.PUERTO_API)+"/login", campos)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al intentar hacer login:", resp.StatusCode)
	}
}

func resetearContraseña(t *testing.T, pass, token string) {
	campos := url.Values{
		"password": {pass},
		"token":    {token},
	}

	resp, err := http.PostForm("http://localhost:"+os.Getenv(globales.PUERTO_API)+"/resetearPassword", campos)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al intentar recibir token:", resp.StatusCode)
	}
}

func recibirTokenResetContraseña(t *testing.T, usuario string) {
	campos := url.Values{
		"usuario": {usuario},
	}

	resp, err := http.PostForm("http://localhost:"+os.Getenv(globales.PUERTO_API)+"/obtenerTokenResetPassword", campos)
	if err != nil {
		t.Fatal("Error al construir request:", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Obtenido código de error no 200 al intentar recibir token:", resp.StatusCode)
	}
}
