package Estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func EjecutarComandoFDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	fmt.Println("############# Comando FDISK ##############")
	ParamValidos = true
	mbr := MBR{}
	particion := Particion{}
	var inicioPart int64 = int64(unsafe.Sizeof(mbr))
	var propiedades [6]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-size":
				propiedades[0] = propiedadTemp.Val
			case "-fit":
				propiedades[1] = propiedadTemp.Val
			case "-unit":
				propiedades[2] = propiedadTemp.Val
			case "-path":
				propiedades[3] = propiedadTemp.Val
			case "-type":
				propiedades[4] = propiedadTemp.Val
			case "-name":
				propiedades[5] = propiedadTemp.Val
				fmt.Println(propiedades[5])

			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}

		//Tamaño de la particion
		var TamanioTotParticioin int64 = 0
		if strings.ToLower(propiedades[2]) == "b" {
			PartitionSize, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotParticioin = PartitionSize
		} else if strings.ToLower(propiedades[2]) == "k" {
			PartitionSize, _ := strconv.ParseInt(propiedades[0], 10, 64)

			TamanioTotParticioin = (PartitionSize * 1024) / 1024

		} else if strings.ToLower(propiedades[2]) == "m" {
			PartitionSize, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotParticioin = (PartitionSize * 1024 * 1024) / 1024
		} else {
			PartitionSize, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotParticioin = (PartitionSize * 1024) / 1024
		}

		//Obtenemos el MBR
		switch strings.ToLower(propiedades[4]) {
		case "p":
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta " + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}
			//Como ya se leyó el mbr, Verificamos si existe espacio disponible o que no lo rebase

			if HayEspacio(TamanioTotParticioin, mbr.Mbr_tamano) {

				return false
			}

			//Verificamos si hay particiones disponibles
			if BytesToString(Particiones[0].Part_status) == "1" {
				// fmt.Println("4")
				// fmt.Println("Ya existe una particon")
				for i := 0; i < 4; i++ {
					inicioPart += Particiones[i].Part_size
					if BytesToString(Particiones[i].Part_status) == "0" {
						// fmt.Println(inicioPart)
						break
					}
				}
			}

			if HayEspacio(inicioPart+TamanioTotParticioin, mbr.Mbr_tamano) {
				return false
			}

			//Aquí le doy valores a la partición
			copy(particion.Part_status[:], "1")
			copy(particion.Part_type[:], propiedades[4])
			copy(particion.Part_fit[:], propiedades[1])
			particion.Part_start = inicioPart
			particion.Part_size = TamanioTotParticioin
			copy(particion.Part_name[:], propiedades[5])
			//Partición creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Part_status) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])

		case "l":
			fmt.Println("Particion Logica")
			if !HayExtendida(propiedades[3]) {
				fmt.Println("No existe una particion Extendida")
				return false
			}
			ebr := EBR{}
			copy(ebr.Part_status[:], "1")
			copy(ebr.Part_fit[:], propiedades[1])
			ebr.Part_start = inicioPart
			ebr.Part_next = 0
			ebr.Part_size = TamanioTotParticioin
			copy(ebr.Part_name[:], propiedades[5])
			//Obtengo el byte donde empieza la partición logica
			IniLogicPart(propiedades[3], ebr)

		case "e":
			//Aquí van las particiones extendidas
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta" + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}
			//Como ya se leyó el mbr, Verificamos si existe espacio disponible o que no lo rebase
			if HayEspacio(TamanioTotParticioin, mbr.Mbr_tamano) {
				return false
			} //Verifico si ya hay particiones
			if BytesToString(Particiones[0].Part_status) == "1" {
				fmt.Println("Ya existe una particion")
				for i := 0; i < 4; i++ {
					//La posición de los bytes del Part_start de la n particion
					inicioPart += Particiones[i].Part_size
					if BytesToString(Particiones[i].Part_status) == "0" {
						fmt.Println(inicioPart)
						break
					}
				}
			}
			if HayEspacio(inicioPart+TamanioTotParticioin, mbr.Mbr_tamano) {
				return false
			}

			//Ahora le damos valores a la partición
			copy(particion.Part_status[:], "1")
			copy(particion.Part_type[:], propiedades[4])
			copy(particion.Part_fit[:], propiedades[1])
			particion.Part_start = inicioPart
			particion.Part_size = TamanioTotParticioin
			copy(particion.Part_name[:], propiedades[5])
			//aquí ya se creó la partición
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Part_status) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
			ebr := EBR{}
			copy(ebr.Part_status[:], "1")
			copy(ebr.Part_fit[:], []byte(propiedades[1]))
			ebr.Part_start = inicioPart
			ebr.Part_next = -1
			ebr.Part_size = TamanioTotParticioin
			copy(ebr.Part_name[:], propiedades[5])
			f.Seek(ebr.Part_start, 0)
			err = binary.Write(f, binary.BigEndian, ebr)
			fmt.Println("Extendida", "Leyendo EBR")
		default:
			fmt.Println("Ocurrió un error")
		}
		return ParamValidos

	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func EliminarParticion(path string, name string, typeDelete string) bool {
	fmt.Print("jsdfl")
	return false
}

func HayEspacio(TamanioTotParticioin int64, tamanioDisco int64) bool {
	if ((TamanioTotParticioin) > tamanioDisco) || (TamanioTotParticioin < 0) {
		fmt.Println("ERROR ----> El tamaño de la particion es mayor a el tamaño del disco o el tamaño es erróneo ")
		return true
	}
	return false
}

func HayExtendida(path string) bool {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			return true
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return false
}

func ReadFile(path string) (funiona bool) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta ")
	}
	fmt.Printf("Fecha: %s\n", mbr.Mbr_fecha_creacion)
	return true
}

func IniLogicPart(path string, ebr2 EBR) bool {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
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
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				ebr.Part_next = ebr.Part_start + int64(unsafe.Sizeof(ebr)) + ebr2.Part_size
				f.Seek(InicioExtendida, 0)
				err = binary.Write(f, binary.BigEndian, ebr)
				ebr2.Part_start = ebr.Part_next
				ebr2.Part_next = -1
				f.Seek(ebr2.Part_start, 0)
				err = binary.Write(f, binary.BigEndian, ebr2)

				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				// fmt.Println(ebr.Part_start)
				// fmt.Println(ebr.Part_next)
				return false
			} else {
				// fmt.Println("Inicio_partición2")
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						fmt.Println("Es la ultima lógica")
						ebr.Part_next = ebr.Part_start + int64(unsafe.Sizeof(ebr)) + ebr2.Part_size
						f.Seek(ebr.Part_start, 0)
						err = binary.Write(f, binary.BigEndian, ebr)
						ebr2.Part_start = ebr.Part_next
						ebr2.Part_next = -1
						f.Seek(ebr2.Part_start, 0)
						err = binary.Write(f, binary.BigEndian, ebr2)
						fmt.Printf("NombreLogica: %s\n", ebr2.Part_name)
						break

					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						fmt.Printf("NombreLogica: %s\n", ebr.Part_name)
					}
				}
				return false
			}

		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return false
}
