package Estructuras

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
)

func EjecutarComandoMKFS(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("############# Comando MKFS ##############")
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-id":
				propiedades[0] = propiedadTemp.Val
			case "-type":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EJecutarMKFS(propiedades[0], ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func EJecutarMKFS(id string, ListaDiscos *list.List) bool {
	dt := time.Now()
	idValido := IdValido(id, ListaDiscos)
	if idValido == false {
		fmt.Println("El id no existe")
		return false
	}
	Id := strings.ReplaceAll(id, "19", "")

	NoParticion := Id[1:]
	IdDisco := Id[1:2]
	pathDisco := ""
	nombreParticion := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disk Disk
		disk = element.Value.(Disk)

		if BytesToString(disk.Id) == IdDisco {
			for i := 0; i < len(disk.Particiones); i++ {
				var mountTemp = disk.Particiones[i]
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")

					nombreParticion = mountTemp.NombreParticion
					pathDisco = disk.Path
					nombreDisco = disk.NombreDisco
					break
				}
			}
		}
		element.Value = disk
	}
	mbr, sizeParticion, InicioParticion := DevolverElMBR(pathDisco, nombreParticion)
	superbloque := SuperBloque{}
	avd := AVD{}
	dd := DD{}
	inodo := INODO{}
	bloque := BLOQUECAR{}
	bitacora := Bitacora{}
	noEstructuras := (sizeParticion - (2 * int64(unsafe.Sizeof(superbloque)))) /
		(27 + int64(unsafe.Sizeof(avd)) + int64(unsafe.Sizeof(dd)) + (5*int64(unsafe.Sizeof(inodo)) +
			(20 * int64(unsafe.Sizeof(bloque))) + int64(unsafe.Sizeof(bitacora))))
	// No. Estructuras

	var cantidadAVD int64 = noEstructuras
	var cantidadDD int64 = noEstructuras
	var cantidadInodos int64 = noEstructuras * 5
	var cantidadBloques int64 = 4 * cantidadInodos
	var Bitacoras int64 = noEstructuras

	//Bitmaaaps
	var InicioBitmapAVD int64 = InicioParticion + int64(unsafe.Sizeof(superbloque))
	var InicioAVD int64 = InicioBitmapAVD + cantidadAVD
	var InicioBitmapDD int64 = InicioAVD + (int64(unsafe.Sizeof(avd)) * cantidadAVD)
	var InicioDD int64 = InicioBitmapDD + cantidadDD
	var InicioBitmapInodo int64 = InicioDD + (int64(unsafe.Sizeof(dd)) * cantidadDD)
	var InicioInodo int64 = InicioBitmapInodo + cantidadInodos
	var InicioBitmapBloque int64 = InicioInodo + (int64(unsafe.Sizeof(inodo)) * cantidadInodos)
	var InicioBLoque int64 = InicioBitmapBloque + cantidadBloques
	var InicioBitacora int64 = InicioBLoque + (int64(unsafe.Sizeof(bloque)) * cantidadBloques)

	//Aquí vamos a inicializar el SuperBloque
	copy(superbloque.Sb_nombre_hd[:], nombreDisco)
	superbloque.Sb_arbol_virtual_count = cantidadAVD
	superbloque.Sb_detalle_directorio_count = cantidadAVD
	superbloque.S_inodes_count = cantidadInodos
	superbloque.S_bloques_count = cantidadBloques

	superbloque.Sb_arbol_virtual_free = cantidadAVD
	superbloque.Sb_detalle_directorio_free = cantidadDD
	superbloque.Sb_inodos_free = cantidadInodos
	superbloque.Sb_bloques_free = cantidadBloques
	copy(superbloque.Sb_date_creacion[:], dt.String())
	copy(superbloque.Sb_date_ultimo_montaje[:], dt.String())
	superbloque.Sb_montajes_count = 1
	//Aquí va todo lo de bitmaps
	superbloque.Sb_ap_bitmap_arbol_directorio = InicioBitmapAVD
	superbloque.Sb_ap_arbol_directorio = InicioAVD
	superbloque.Sb_ap_bitmap_detalle_directorio = InicioBitmapDD
	superbloque.Sb_ap_detalle_directorio = InicioDD
	superbloque.Sb_ap_bitmap_tabla_inodo = InicioBitmapInodo
	superbloque.Sb_ap_tabla_inodo = InicioInodo
	superbloque.Sb_ap_bitmap_bloques = InicioBitmapBloque
	superbloque.Sb_ap_bloques = InicioBLoque
	superbloque.Sb_ap_log = InicioBitacora
	superbloque.Sb_size_struct_arbol_directorio = int64(unsafe.Sizeof(avd))
	superbloque.Sb_size_struct_Detalle_directorio = int64(unsafe.Sizeof(dd))
	superbloque.Sb_size_struct_inodo = int64(unsafe.Sizeof(inodo))
	superbloque.Sb_size_struct_bloque = int64(unsafe.Sizeof(bloque))
	superbloque.Sb_first_free_bit_arbol_directorio = InicioBitmapAVD
	superbloque.Sb_first_free_bit_detalle_directorio = InicioBitmapDD
	superbloque.Sb_dirst_free_bit_tabla_inodo = InicioBitmapInodo
	superbloque.Sb_first_free_bit_bloques = InicioBitmapBloque
	superbloque.Sb_magic_num = 123
	superbloque.ConteoAVD = 0
	superbloque.ConteoDD = 0
	superbloque.ConteoInodo = 0
	superbloque.ConteoBloque = 0
	//Escribo en la particion
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()

	//Escribo el Super Boot
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &superbloque)
	//Escribo el bit map Arbol virtual de directorio
	f.Seek(InicioBitmapAVD, 0)
	var otro int8 = 0
	var i int64 = 0
	for i = 0; i < cantidadAVD; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribo el arbol de directorio
	f.Seek(InicioAVD, 0)
	i = 0
	for i = 0; i < cantidadAVD; i++ {
		err = binary.Write(f, binary.BigEndian, &avd)
	}
	//Escribir el bitmap del detalle directorio
	f.Seek(InicioBitmapDD, 0)
	i = 0
	for i = 0; i < cantidadDD; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribo el detalle en el directorio
	f.Seek(InicioDD, 0)
	i = 0
	dd.Dd_ap_detalle_directorio = -1
	for i = 0; i < cantidadDD; i++ {
		err = binary.Write(f, binary.BigEndian, &dd)
	}
	//Escribir Bitmap en la tabla Inodo
	f.Seek(InicioBitmapInodo, 0)
	i = 0
	for i = 0; i < cantidadInodos; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribir en la tabal Inodos
	f.Seek(InicioInodo, 0)
	i = 0
	inodo.I_count_inodo = -1
	for i = 0; i < cantidadInodos; i++ {
		err = binary.Write(f, binary.BigEndian, &inodo)
	}

	//Escribo el bitmap bloque de datos
	f.Seek(InicioBitmapBloque, 0)
	i = 0
	for i = 0; i < cantidadBloques; i++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	//Escribo bloque de datos
	f.Seek(InicioBLoque, 0)
	i = 0
	copy(bloque.B_content[:], "")
	for i = 0; i < cantidadBloques; i++ {
		err = binary.Write(f, binary.BigEndian, &bloque)
	}
	//Escribir Bitacoras
	f.Seek(InicioBitacora, 0)
	i = 0
	bitacora.Size = -1
	for i = 0; i < Bitacoras; i++ {
		err = binary.Write(f, binary.BigEndian, &bitacora)
	}

	CreateRoot(pathDisco, InicioParticion)
	// fmt.Println("No. Estructuras:", noEstructuras)
	fmt.Println("Particion a formatear", nombreParticion, NoParticion)
	fmt.Println("Tamaño de la particion", sizeParticion)
	fmt.Println("Fecha: %s\n", mbr.Mbr_fecha_creacion)
	return false
}

func DevolverElMBR(path string, nombreParticion string) (MBR, int64, int64) {
	mbr := MBR{}
	var Particiones [4]Particion
	var nombre2 [15]byte
	var size int64
	copy(nombre2[:], nombreParticion)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return mbr, 0, 0
	}
	defer f.Close()
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	Particiones = mbr.Particiones
	for i := 0; i < 4; i++ {
		if BytesNombreParticion(Particiones[i].Part_name) == BytesNombreParticion(nombre2) {
			size = Particiones[i].Part_size
			return mbr, size, Particiones[i].Part_start
		}
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				fmt.Println("No hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					if BytesNombreParticion(ebr.Part_name) == BytesNombreParticion(nombre2) {
						fmt.Println("Logica Encontrada")
						return mbr, ebr.Part_size, ebr.Part_start
					}
				}
			}
		}
	}
	return mbr, 0, 0
}

func CreateRoot(pathDisco string, InicioParticion int64) bool {
	dt := time.Now()
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(InicioParticion, 0)
	sb := SuperBloque{}
	err = binary.Read(f, binary.BigEndian, &sb)
	//Escribo 1 en el bitmap avd y escribo en el avd

	f.Seek(sb.Sb_ap_bitmap_arbol_directorio, 0)
	var otro int8 = 0
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	bitLibre, _ := f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_arbol_directorio = bitLibre
	avd := AVD{}
	copy(avd.Avd_fecha_creacion[:], dt.String())
	copy(avd.Avd_nombre_directotrio[:], "/")
	for j := 0; j < 6; j++ {
		avd.Avd_ap_array_subdirectorios[j] = -1
	}
	avd.Avd_ap_detalle_directorio = 0
	avd.Avd_ap_arbol_virtual_directorio = -1
	copy(avd.Avd_proper[:], "root")
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	err = binary.Write(f, binary.BigEndian, &avd)

	sb.Sb_arbol_virtual_free = sb.Sb_arbol_virtual_free - 1
	//Escribo 1 en el bitmap detalle directori y escribo en el detalledirectorio
	f.Seek(sb.Sb_ap_bitmap_detalle_directorio, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_detalle_directorio = bitLibre
	detalleDirectorio := DD{}
	arregloDD := ArregloDD{}
	copy(arregloDD.Dd_file_nombre[:], "user.txt")
	copy(arregloDD.Dd_file_date_creacion[:], dt.String())
	copy(arregloDD.Dd_file_date_modificacion[:], dt.String())
	arregloDD.Dd_file_ap_inodo = 0
	detalleDirectorio.Dd_array_files[0] = arregloDD
	detalleDirectorio.Ocupado = 1
	for j := 0; j < 5; j++ {
		if j == 0 {
			detalleDirectorio.Dd_array_files[j].Dd_file_ap_inodo = 0
		} else {
			detalleDirectorio.Dd_array_files[j].Dd_file_ap_inodo = -1
		}
	}
	detalleDirectorio.Dd_ap_detalle_directorio = -1
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	err = binary.Write(f, binary.BigEndian, &detalleDirectorio)

	sb.Sb_detalle_directorio_free = sb.Sb_detalle_directorio_free - 1

	// Escribo 1 en el bitmap tablaInodo y escribo en el inodo

	var cantidadBloque int64 = CantidadBloqueUsar("1,G,root\n1,U,root,root,123\n")
	f.Seek(sb.Sb_ap_bitmap_tabla_inodo, 0)
	otro = 1
	err = binary.Write(f, binary.BigEndian, &otro)
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_dirst_free_bit_tabla_inodo = bitLibre
	inodo := INODO{}
	for j := 0; j < 4; j++ {
		inodo.I_array_bloques[j] = -1
	}
	inodo.I_count_inodo = 0
	inodo.I_size_archivo = 10
	inodo.I_count_bloques_asignados = cantidadBloque
	for h := 0; h < int(cantidadBloque); h++ {
		inodo.I_array_bloques[h] = int64(h)
	}
	inodo.I_ao_indirecto = -1
	inodo.I_id_proper = 123
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	err = binary.Write(f, binary.BigEndian, &inodo)
	sb.Sb_inodos_free = sb.Sb_inodos_free - 1
	//Escribo 1 en el bitmap BloqueDatos y escribo en el bloque de datos

	f.Seek(sb.Sb_ap_bitmap_bloques, 0)
	otro = 1
	for k := 0; k < int(cantidadBloque); k++ {
		err = binary.Write(f, binary.BigEndian, &otro)
	}
	otro = 0
	bitLibre, _ = f.Seek(0, os.SEEK_CUR)
	sb.Sb_first_free_bit_bloques = bitLibre
	f.Seek(sb.Sb_ap_bloques, 0)
	usesTxt := []byte("1,G,root\n1,U,root,root,123\n")
	for k := 0; k < int(cantidadBloque); k++ {
		if k == 0 {
			bloque := BLOQUECAR{}
			copy(bloque.B_content[:], string([]byte(usesTxt[0:25])))
			err = binary.Write(f, binary.BigEndian, &bloque)
		} else {
			bloque := BLOQUECAR{}
			copy(bloque.B_content[:], string([]byte(usesTxt[k*25:len(usesTxt)])))
			err = binary.Write(f, binary.BigEndian, &bloque)
		}
		sb.Sb_bloques_free = sb.Sb_bloques_free - 1
		sb.ConteoBloque = sb.ConteoBloque + int64(k)
	}

	//Actualizo el superbloque

	f.Seek(0, 0)
	f.Seek(InicioParticion, 0)
	err = binary.Write(f, binary.BigEndian, &sb)
	return false

}

func CantidadBloqueUsar(data string) int64 {
	var noBloque int64 = 0
	cont := 1
	var dataX []byte = []byte(data)
	for i := 0; i < len(dataX); i++ {
		if cont == 25 {
			noBloque = noBloque + 1
			cont = 0
		}
		cont++
	}
	if len(dataX)%25 != 0 {
		noBloque = noBloque + 1
	}
	return noBloque

}

func DevolverSuperBloque(path string, nombreParticion string) (SuperBloque, int64) {
	mbr := MBR{}
	sb := SuperBloque{}
	var Particiones [4]Particion
	var nombre2 [15]byte
	copy(nombre2[:], nombreParticion)
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return sb, 0
	}
	defer f.Close()

	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	Particiones = mbr.Particiones
	for i := 0; i < 4; i++ {
		if BytesNombreParticion(Particiones[i].Part_name) == BytesNombreParticion(nombre2) {
			f.Seek(Particiones[i].Part_start, 0)
			err = binary.Read(f, binary.BigEndian, &sb)
			return sb, Particiones[i].Part_start
		}
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					if BytesNombreParticion(ebr.Part_name) == BytesNombreParticion(nombre2) {
						fmt.Println("Logica Encontrada")
						f.Seek(ebr.Part_start, 0)
						err = binary.Read(f, binary.BigEndian, &sb)
						return sb, ebr.Part_start
					}

				}
			}
		}
	}
	return sb, 0
}
