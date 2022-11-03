package Estructuras

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func EjecutarComandoMKFILE(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("############# Comando MKFILE #############")
	ParamValidos = true
	pathEspacio := " "
	var propiedades [5]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-id":
				propiedades[0] = propiedadTemp.Val
			case "-path":
				propiedades[1] = propiedadTemp.Val
			case "-r":
				propiedades[2] = propiedadTemp.Val
			case "-size":
				propiedades[3] = propiedadTemp.Val
			case "-cont":
				propiedades[4] = propiedadTemp.Val
			case "-sigue":
				propiedades[4] += propiedadTemp.Val
			case "-sigueCont":
				fmt.Println(propiedadTemp.Val + " ")
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		size, _ := strconv.Atoi(propiedades[3])

		if pathEspacio != " " {
			pathCompleta := propiedades[1] + pathEspacio
			EjecutarMKFILE(global2, pathCompleta[1:len(pathCompleta)-2], propiedades[2], size, propiedades[4], ListaDiscos)
		} else {
			EjecutarMKFILE(global2, propiedades[1], propiedades[2], size, propiedades[4][0:len(propiedades[4])-1], ListaDiscos)
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func EjecutarMKFILE(idParticion string, pathArchivo string, _p string, size int, contenido string, ListaDiscos *list.List) bool {
	if size > len(contenido) {
		for i := len(contenido); i < size; i++ {
			contenido = contenido + " "
		}
	}

	//Aquí quito las comillas al path si es que tiene

	Escomilla := pathArchivo[0:1]
	if Escomilla == "\"" {
		pathArchivo = pathArchivo[1 : len(pathArchivo)-1]
	}
	pathDisco, nombreParticion, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	CreateFile(pathDisco, nombreParticion, pathArchivo, _p, size, contenido, -1)
	return true
}

func CreateFile(pathDisco string, nombreParticion string, pathArchivo string, _p string, size int, contenido string, siguienteDD int) bool {
	// Obtener el SB de la particion

	otroDD := true
	dt := time.Now()
	avd := AVD{}
	superbloque := SuperBloque{}
	encontrado := false
	dd := DD{}

	var InicioParticion int64 = 0
	var nombreArchivo = ""
	var carpetaPadre = ""

	if strings.Contains(pathArchivo, "/") {
		nuevaPath := "/"
		carpetas := strings.Split(pathArchivo, "/")
		nombreArchivo = carpetas[len(carpetas)-1]
		carpetaPadre = carpetas[len(carpetas)-2]
		for i := 1; i < len(carpetas)-1; i++ {
			nuevaPath += carpetas[i] + "/"
		}
		nuevaPath = nuevaPath[0 : len(nuevaPath)-1]
		//Se crean las carpetas si no estan creadas
		RecorrerElPath(nuevaPath, nombreParticion, pathDisco)
	}
	superbloque, InicioParticion = DevolverSuperBloque(pathDisco, nombreParticion)
	//Primero busco el AVD, Luego el DD, Luego modifico el DD, luego creo el Inodo, luego creo el Bloque, luego modifico AVD
	//Obtenemos el AVD
	var nombre2 [15]byte
	var bitLibre int64 = 0
	copy(nombre2[:], carpetaPadre)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()

	//Escribo en la bitacora
	if siguienteDD == -1 {
		f.Seek(superbloque.Sb_ap_log, 0)
		bitacora := Bitacora{}
		copy(bitacora.Log_tipo_operacion[:], "mkfile")
		copy(bitacora.Log_tipo[:], "1")
		copy(bitacora.Log_nombre[:], pathArchivo)
		copy(bitacora.Log_Contenido[:], contenido[1:len(contenido)-1])
		copy(bitacora.Log_fecha[:], dt.String())
		bitacora.Size = int64(size)
		bitacoraTemp := Bitacora{}
		var bitBitacora int64 = 0
		for i := 0; i < 3000; i++ {
			bitBitacora, _ = f.Seek(0, os.SEEK_CUR)
			err = binary.Read(f, binary.BigEndian, &bitacoraTemp)
			if bitacoraTemp.Size == -1 {
				f.Seek(bitBitacora, 0)
				err = binary.Write(f, binary.BigEndian, &bitacora)
				break
			}
		}
	}
	//Escribo en el Arbol Directorio
	f.Seek(superbloque.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(superbloque.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if BytesNombreParticion(avd.Avd_nombre_directotrio) == BytesNombreParticion(nombre2) {
			//Leo el DD
			f.Seek(superbloque.Sb_ap_detalle_directorio, 0)
			for i := 0; i < 20; i++ {
				err = binary.Read(f, binary.BigEndian, &dd)
				if dd.Dd_ap_detalle_directorio != -1 && dd.Dd_ap_detalle_directorio != 0 {
					siguienteDD = int(dd.Dd_ap_detalle_directorio)
					bitLibre, _ = f.Seek(0, os.SEEK_CUR)
					continue
				}
				if i == int(avd.Avd_ap_detalle_directorio) || i == siguienteDD {
					for i := 0; i < len(dd.Dd_array_files); i++ {
						if dd.Dd_array_files[i].Dd_file_ap_inodo == -1 {
							copy(dd.Dd_array_files[i].Dd_file_nombre[:], nombreArchivo)
							dd.Dd_array_files[i].Dd_file_ap_inodo = superbloque.ConteoInodo + 1
							copy(dd.Dd_array_files[i].Dd_file_date_creacion[:], dt.String())
							copy(dd.Dd_array_files[i].Dd_file_date_modificacion[:], dt.String())
							otroDD = false
							break
						}
					}
					if otroDD == false {
						f.Seek(bitLibre, 0)
						err = binary.Write(f, binary.BigEndian, &dd)
						bitLibre = 0
						encontrado = true
						WriteInodo(pathDisco, superbloque, contenido[1:len(contenido)-1], InicioParticion)
						break
					} else if otroDD == true {
						f.Seek(bitLibre, 0)
						dd.Dd_ap_detalle_directorio = superbloque.ConteoDD + 1
						err = binary.Write(f, binary.BigEndian, &dd)
						bitLibre = 0
						f.Seek(0, 0)
						//Creo otro detalle de directorio
						nuevoDD := DD{}
						superbloque.ConteoDD = superbloque.ConteoDD + 1
						nuevoDD.Ocupado = 1
						//Marco 1 bitmap DD
						superbloque.Sb_detalle_directorio_free = superbloque.Sb_detalle_directorio_free - 1
						f.Seek(superbloque.Sb_first_free_bit_detalle_directorio, 0)
						var otro int8 = 1
						err = binary.Write(f, binary.BigEndian, &otro)
						otro = 0
						bitLibre, _ = f.Seek(0, os.SEEK_CUR)
						superbloque.Sb_first_free_bit_detalle_directorio = bitLibre
						//Actualizo mi superbloque
						f.Seek(InicioParticion, 0)
						err = binary.Write(f, binary.BigEndian, &superbloque)
						WriteDD(superbloque.Sb_ap_detalle_directorio, pathDisco, superbloque.Sb_detalle_directorio_count, nuevoDD)
						f.Seek(0, 0)
						CreateFile(pathDisco, nombreParticion, pathArchivo, _p, size, contenido, int(dd.Dd_ap_detalle_directorio))
						encontrado = true
						break

					}
				}
				bitLibre, _ = f.Seek(0, os.SEEK_CUR)
			}
		}
		if encontrado == true {
			break
		}

	}
	return false
}

func WriteDD(InicioDD int64, pathDisco string, cantidadDD int64, ddNuevo DD) bool {
	dd := DD{}
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(InicioDD, 0)
	var bitLibre int64
	for i := 0; i < int(cantidadDD); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			f.Seek(bitLibre, 0)
			for h := 0; h < 5; h++ {
				ddNuevo.Dd_array_files[h].Dd_file_ap_inodo = -1
			}
			ddNuevo.Dd_ap_detalle_directorio = -1
			err = binary.Write(f, binary.BigEndian, &ddNuevo)
			break
		}
		bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	}
	return false
}

func WriteInodo(pathDisco string, superbloque SuperBloque, contenido string, InicioParticion int64) bool {
	var otro int8 = 0
	var bitLibre int64 = 0
	var restoBloque int64 = 0
	contenido2 := ""
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	//Escribo en el bitmap de Inodo
	var cantidadBloque int64 = CantidadBloqueUsar(contenido)
	if cantidadBloque >= 5 {
		restoBloque = cantidadBloque - 4
		cantidadBloque = 4
		contenido2 = contenido[100:len(contenido)]
		contenido = contenido[0:100]
	}
	f.Seek(superbloque.Sb_dirst_free_bit_tabla_inodo, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	superbloque.Sb_dirst_free_bit_tabla_inodo = bitLibre
	inodo := INODO{}
	for j := 0; j < 4; j++ {
		inodo.I_array_bloques[j] = -1
	}
	inodo.I_count_inodo = superbloque.ConteoInodo + 1
	inodo.I_size_archivo = 10
	inodo.I_count_bloques_asignados = cantidadBloque
	inodo.I_ao_indirecto = -1
	inodo.I_id_proper = 123

	inodo, superbloque.ConteoBloque, superbloque.Sb_first_free_bit_bloques, superbloque.Sb_bloques_free = WriteBloque(superbloque, cantidadBloque, pathDisco, InicioParticion, inodo, contenido)
	f.Seek(superbloque.Sb_ap_tabla_inodo, 0)
	inodoTemp := INODO{}
	superbloque.ConteoInodo = superbloque.ConteoInodo + 1
	superbloque.Sb_inodos_free = superbloque.Sb_inodos_free - 1
	if restoBloque != 0 {
		inodo.I_ao_indirecto = superbloque.ConteoInodo + 1
	}
	for i := 0; i < int(superbloque.S_inodes_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodoTemp)
		if inodoTemp.I_count_inodo == -1 {
			f.Seek(bitLibre, 0)
			err = binary.Write(f, binary.BigEndian, &inodo)
			break
		}
		bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	}
	//Actualizo mi superbloque
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &superbloque)
	if restoBloque != 0 {
		WriteInodo(pathDisco, superbloque, contenido2, InicioParticion)
	}
	return false
}

func WriteBloque(superbloque SuperBloque, cantidadBloque int64, pathDisco string, InicioParticion int64, inodo INODO, contenido string) (INODO, int64, int64, int64) {
	var contenido2 [25]byte
	copy(contenido2[:], contenido)
	bloqueTemp := BLOQUECAR{}
	var bitLibre_Bloque int64 = 0
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return inodo, 0, 0, 0
	}
	defer f.Close()
	//Escribo en el bitmap en el bloque de datos
	var otro int8 = 0
	f.Seek(superbloque.Sb_first_free_bit_bloques, 0)
	otro = 1
	for k := 0; k < int(cantidadBloque); k++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	otro = 0
	bitLibre, _ := f.Seek(0, os.SEEK_CUR)
	superbloque.Sb_first_free_bit_bloques = bitLibre
	//Escribo los bloques de datos

	f.Seek(superbloque.Sb_ap_bloques, 0)
	for i := 0; i < int(superbloque.S_bloques_count); i++ {
		err = binary.Read(f, binary.BigEndian, &bloqueTemp)
		if bloqueTemp.B_content[0] == 0 {
			f.Seek(bitLibre_Bloque, 0)
			for h := 0; h < int(cantidadBloque); h++ {
				inodo.I_array_bloques[h] = superbloque.ConteoBloque + 1
				//Escribo el bloque
				if h == 0 {
					bloque := BLOQUECAR{}
					if len(contenido) >= 25 {
						copy(bloque.B_content[:], string([]byte(contenido[0:25])))
						err = binary.Write(f, binary.BigEndian, &bloque)
					} else {
						bloque.B_content = contenido2
						err = binary.Write(f, binary.BigEndian, &bloque)
					}

				} else {
					bloque := BLOQUECAR{}
					copy(bloque.B_content[:], string([]byte(contenido[h*25:len(contenido)])))
					err = binary.Write(f, binary.BigEndian, &bloque)
				}
				superbloque.Sb_bloques_free = superbloque.Sb_bloques_free - 1
				superbloque.S_bloques_count = superbloque.S_bloques_count + 1
			}
			break
		}
		bitLibre_Bloque, _ = f.Seek(0, os.SEEK_CUR)
	}
	return inodo, superbloque.ConteoBloque, superbloque.Sb_first_free_bit_bloques, superbloque.Sb_bloques_free
}
