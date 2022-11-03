package Estructuras

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func EjecutarComandoMKGRP(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("############# Comando MKGRP ##############")
	ParamValidos = true
	var propiedades [2]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-id":
				propiedades[0] = propiedadTemp.Val
			case "-name":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EjecutarMKGRP(global2, propiedades[1], ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func EjecutarMKGRP(id string, name string, ListaDiscos *list.List) {
	pathDisco, nombreParticion, _ := RecorrerListaDisco(id, ListaDiscos)
	if global == "root" {
		ModificateFile(pathDisco, nombreParticion, "user.txt", name)
	} else {
		fmt.Println("No hay usuario logeado / no es usuario root")
	}
}

func ModificateFile(pathDisco string, nombreParticion string, nombreArchivo string, group string) bool {
	superbloque := SuperBloque{}
	superbloque, _ = DevolverSuperBloque(pathDisco, nombreParticion)
	avd := AVD{}
	dd := DD{}
	inodo := INODO{}
	f, err := os.OpenFile(pathDisco, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(superbloque.Sb_ap_arbol_directorio, 0)
	err = binary.Read(f, binary.BigEndian, &avd)

	//AVD ya está inicializado
	apuntadorDD := avd.Avd_ap_detalle_directorio
	f.Seek(superbloque.Sb_ap_detalle_directorio, 0)
	for i := 0; i < int(superbloque.Sb_detalle_directorio_count); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if i == int(apuntadorDD) {
			break
		}
	}
	arregloDD := ArregloDD{}
	arregloDD = dd.Dd_array_files[0]
	apuntadorInodo := arregloDD.Dd_file_ap_inodo
	f.Seek(superbloque.Sb_ap_tabla_inodo, 0)
	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if i == int(apuntadorInodo) {
			break
		}
	}
	fmt.Printf("Archivo: %s\n", arregloDD.Dd_file_nombre)
	return false

}
