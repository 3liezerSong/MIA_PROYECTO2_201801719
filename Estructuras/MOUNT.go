package Estructuras

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func EjejcutarComandoMount(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("############# Comando MOUNT ##############")
	var propiedades [2]string
	var nombre [15]byte
	ParamValidos = true
	if len(propiedadesTemp) >= 2 {
		//Recorro la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name

			switch strings.ToLower(nombrePropiedad) {
			case "-name":
				propiedades[0] = propiedadTemp.Val
				copy(nombre[:], propiedades[0])
			case "-path":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		//Empiezo a montar las particiones
		EjecutCommand(propiedades[1], nombre, ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func EjecutCommand(path string, NombreParticion [15]byte, ListaDiscos *list.List) bool {
	var encontrada = false
	lineaComando := strings.Split(path, "/")
	nombreDisco := lineaComando[len(lineaComando)-1]
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if string(Particiones[i].Part_name[:]) == string(NombreParticion[:]) {
			encontrada = true
			if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
				fmt.Println("Error no es posible montar una partición Extendida")

			} else {
				MontarParticion(ListaDiscos, string(NombreParticion[:]), string(nombreDisco), path)
			}
		}
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			ebr := EBR{}
			f.Seek(Particiones[i].Part_start, 0)
			err = binary.Read(f, binary.BigEndian, &ebr)
			for {
				if ebr.Part_next == -1 {
					break
				} else {
					f.Seek(ebr.Part_next, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
				}
				var nombre string = string(ebr.Part_name[:])
				var nombre2 string = string(NombreParticion[:])
				if nombre == nombre2 {

					encontrada = true
					//Montar Partition
					MontarParticion(ListaDiscos, string(NombreParticion[:]), string(nombreDisco), path)
				}
			}
		}
	}
	if encontrada == false {
		fmt.Println("Error no se encontró la partición")
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta especificada")
	}
	return true
}

func MontarParticion(ListaDiscos *list.List, nombreParticion string, nombreDisco string, path string) {
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disk Disk
		disk = element.Value.(Disk)

		if BytesToString(disk.Estado) == "0" && !DiskExist(ListaDiscos, nombreDisco) {
			disk.NombreDisco = nombreDisco
			disk.Path = path
			copy(disk.Estado[:], "1")
			//-id -> 191A
			for i := 0; i < len(disk.Particiones); i++ {
				var mountTemp = disk.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" { // En esta parte voy generando los id para los mount
					mountTemp.Id = "19" + mountTemp.Id + BytesToString(disk.Id)
					mountTemp.NombreParticion = nombreParticion
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disk.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.NombreParticion == nombreParticion {
					fmt.Println("La particion ya está montada")
					break
				}
			}
			element.Value = disk
			break
		} else if BytesToString(disk.Estado) == "1" && DiskExist(ListaDiscos, nombreDisco) && nombreDisco == disk.NombreDisco {

			fmt.Println("Otra particion montada en el disco ", BytesToString(disk.Id)+" Aumentamos el numero")
			for i := 0; i < len(disk.Particiones); i++ {
				var mountTemp = disk.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "19" + mountTemp.Id + BytesToString(disk.Id)
					mountTemp.NombreParticion = nombreParticion
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disk.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.NombreParticion == nombreParticion {
					fmt.Println("La partición ya se encuentra montada")
					break
				}
			}
			element.Value = disk
			break
		}

	}
}

func DiskExist(ListaDiscos *list.List, nombreDisco string) bool {
	Exist := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disk Disk
		disk = element.Value.(Disk)
		if disk.NombreDisco == nombreDisco {

			return true
		} else {

			Exist = false
		}
	}
	return Exist
}

func IdValido(id string, ListaDiscos *list.List) bool {
	esta := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disk Disk
		disk = element.Value.(Disk)
		if disk.NombreDisco != "" {

			for i := 0; i < len(disk.Particiones); i++ {
				var mountTemp = disk.Particiones[i]
				if mountTemp.NombreParticion != "" {
					if mountTemp.Id == id {
						return true
					}
				}
			}
		}
	}
	return esta
}
