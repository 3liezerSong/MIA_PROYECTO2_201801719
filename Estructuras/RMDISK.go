package Estructuras

import (
	"fmt"
	"os"
	"strings"
)

func EjecutarComandoRMDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	fmt.Println("############# Comando RMDISK #############")
	ParamValidos = true
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-path":
				err := os.Remove(propiedadTemp.Val)
				if err != nil {
					fmt.Println("Error, no se encontró el disco")
				} else {
					fmt.Println("Disco Elimindo Correctamente")
				}

			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
