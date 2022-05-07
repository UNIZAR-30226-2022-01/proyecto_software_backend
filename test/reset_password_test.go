package integracion

import (
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/dao"
	"github.com/UNIZAR-30226-2022-01/proyecto_software_backend/globales"
	"testing"
	"time"
)

func TestResetContraseña(t *testing.T) {
	// Se debe ejecutar manualmente, ya que se debe rellenar el email de destino y asegurar que el fichero ../mail.env existe
	t.Skip()
	email := "" // Rellenar con dirección email accesible

	t.Log("Purgando DB...")
	purgarDB()

	// Creación e inicio de la partida

	t.Log("Creando usuarios...")
	_ = crearUsuarioConEmail("usuario1", email, t)

	// Pide un reseteo
	recibirTokenResetContraseña(t, "usuario1")

	// Comprueba el token
	token := dao.ObtenerToken(globales.Db, "usuario1")

	t.Log("token:", token)

	// Resetea la contraseña e intenta loguear con la introducida
	resetearContraseña(t, "nuevaPass", token)
	login(t, "usuario1", "nuevaPass")

	time.Sleep(120 * time.Second)
}
