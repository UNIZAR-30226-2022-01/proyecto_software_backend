package main

import (
	"backend/globales"
	"backend/servidor"
)

func main() {
	globales.InicializarGrafoMapa()

	servidor.IniciarServidor(false)
}
