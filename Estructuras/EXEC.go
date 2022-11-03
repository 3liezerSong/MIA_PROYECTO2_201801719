package Estructuras

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"strings"
)

func Check(e error) {
	if e != nil {
		fmt.Println("Error")
	}
}
func EjecutarComandoExec(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("############# Comando EXEC ##############")
	ParamValidos = true
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-path":
				dat, err := ioutil.ReadFile(strings.Replace(propiedadTemp.Val, "\"", "", 2))
				Check(err)
				AnalizarComando(string(dat), ListaDiscos)
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
