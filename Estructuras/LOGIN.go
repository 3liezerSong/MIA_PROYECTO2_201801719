package Estructuras

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func EjecutarComandoLogin(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (bool, string, string) {
	fmt.Println("############# Comando Login ##############")
	ParamValidos := true
	usuario := ""
	iddisk := ""
	var propiedades [3]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-usuario":
				propiedades[0] = propiedadTemp.Val
			case "-password":
				propiedades[1] = string(propiedadTemp.Val)
			case "-id":
				propiedades[2] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el comando")
			}
		}
		ParamValidos, usuario, iddisk = EjecutarLogin(propiedades[0], propiedades[1], propiedades[2], ListaDiscos)
		return ParamValidos, usuario, iddisk

	} else {
		ParamValidos = false
		return ParamValidos, usuario, iddisk

	}
}

func EjecutarLogin(usuario string, password string, id string, ListaDiscos *list.List) (bool, string, string) {
	idValido := IdValido(id, ListaDiscos)

	if idValido == false {
		fmt.Println("El id no existe o puecde ser que la particion no está montada")
		return false, "", ""
	} else if global != "" {
		fmt.Println("Ya hay una sesion iniciada")
		return false, "", ""
	}
	pathDisco, nombreParticon, nombreDisco := RecorrerListaDisco(id, ListaDiscos)
	mbr, sizeParticion, InicioParticion := DevolverElMBR(pathDisco, nombreParticon)
	superBloque := SuperBloque{}
	f, err := os.OpenFile(pathDisco, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false, "", ""
	}

	defer f.Close()
	f.Seek(InicioParticion, 0)
	err = binary.Read(f, binary.BigEndian, &superBloque)

	//Obtengo el avd root

	avd := AVD{}
	dd := DD{} //detalle del directorio
	inodo := INODO{}
	bloque := BLOQUECAR{}
	f.Seek(superBloque.Sb_ap_arbol_directorio, 0)
	err = binary.Read(f, binary.BigEndian, &avd)
	apuntadorDD := avd.Avd_ap_detalle_directorio
	f.Seek(superBloque.Sb_ap_detalle_directorio, 0)
	for i := 0; i < int(superBloque.Sb_arbol_virtual_free); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if i == int(apuntadorDD) {
			break
		}
	}

	apuntadorInodo := dd.Dd_array_files[0].Dd_file_ap_inodo
	f.Seek(superBloque.Sb_ap_tabla_inodo, 0)
	for i := 0; i < int(superBloque.Sb_inodos_free); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if i == int(apuntadorInodo) {
			break
		}
	}
	var usertxt string = ""

	//Aqui voy a leer el archivo Users.txt
	posicion := 0
	f.Seek(superBloque.Sb_ap_bloques, 0)

	for i := 0; i < int(superBloque.Sb_inodos_free); i++ {

		err = binary.Read(f, binary.BigEndian, &bloque)

		if int(inodo.I_array_bloques[posicion]) != -1 && int(inodo.I_array_bloques[posicion]) == i {

			usertxt += ConvertirData(bloque.B_content)
		} else if int(inodo.I_array_bloques[posicion]) == -1 {

			break
		} else {
			break
		}
		if posicion < 4 {
			posicion++
		} else if posicion == 4 {
			posicion = 0
		}

	}

	lineUsuTxt := strings.Split(usertxt, "\n")
	for i := 0; i < len(lineUsuTxt)-1; i++ {
		if len(lineUsuTxt[i]) != 19 {
			groupusu := strings.Split(lineUsuTxt[i], ",")
			fmt.Println(groupusu)
			if groupusu[1] == "U" {
				if groupusu[3] == usuario && groupusu[4] == password {
					fmt.Println("Bienvendio al sistema")
					return true, usuario, id
				}
			}
		}
	}
	fmt.Println(nombreDisco, mbr.Mbr_tamano, sizeParticion)
	return false, "", id

}
