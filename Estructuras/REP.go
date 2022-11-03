package Estructuras

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func EjecutarComandoRep(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("############# Haciendo Reporte #############")
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
			case "-path":
				propiedades[1] = propiedadTemp.Val
			case "-name":
				propiedades[2] = propiedadTemp.Val
			case "-ruta":
				propiedades[3] = propiedadTemp.Val
			case "-sigue":
				propiedades[1] += propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando", nombrePropiedad)
			}
		}
		EsComilla := propiedades[1][0:1]
		if EsComilla == "\"" {
			if propiedades[3] != "" {
				propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
			}
			propiedades[1] = propiedades[1][1 : len(propiedades[1])-1]
		}
		carpetas_Graficar := strings.Split(propiedades[1], "/")
		var comando = ""
		for i := 1; i < len(carpetas_Graficar)-1; i++ {
			comando += carpetas_Graficar[i] + "/"
		}
		fmt.Println(comando)
		executeComand("mkdir " + comando[0:len(comando)-1])
		switch strings.ToLower(propiedades[2]) {
		case "disk":
			GraficarDisco(propiedades[0], ListaDiscos, propiedades[1])

		case "file":
			Reportefile(propiedades[0], propiedades[1], propiedades[3], ListaDiscos)
		case "tree":
			GraficarTree(propiedades[0], propiedades[1], propiedades[3], ListaDiscos)
		default:
			fmt.Println("Nombre Incorrecto")

		}
		return ParamValidos

	} else {
		ParamValidos = false
		return ParamValidos
	}
}

//Graficar Disco y calcular Porcentajes
func GraficarDisco(idParticion string, ListaDiscos *list.List, path string) bool {
	var NombreParticion [15]byte
	var buffer bytes.Buffer
	buffer.WriteString("digraph G{\ntbl [\nshape=box\nlabel=<\n<table border='0' cellborder='2' width='100' height=\"30\" color='lightblue4'>\n<tr>")
	pathDisco, _, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco + "Jajajj")
		return false
	}
	defer f.Close()
	PorcentajeUtilizado := 0.0
	var EspacioUtilizado int64 = 0
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	TamanioDisco := mbr.Mbr_tamano
	Particiones := mbr.Particiones
	buffer.WriteString("<td height='30' width='75'> MBR </td>")
	for i := 0; i < 4; i++ {
		if convertName(Particiones[i].Part_name[:]) != convertName(NombreParticion[:]) && strings.ToLower(BytesToString(Particiones[i].Part_type)) == "p" {
			PorcentajeUtilizado = (float64(Particiones[i].Part_size) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td height='30' width='75.0'>PRIMARIA <br/>" + convertName(Particiones[i].Part_name[:]) + " <br/> Ocupado: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
			EspacioUtilizado += Particiones[i].Part_size
		} else if convertName(Particiones[i].Part_status[:]) == "0" {
			buffer.WriteString("<td height='30' width='75.0'>Libre</td>")
		}
		if strings.ToLower(BytesToString(Particiones[i].Part_type)) == "e" {
			EspacioUtilizado += Particiones[i].Part_size
			PorcentajeUtilizado = (float64(Particiones[i].Part_size) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td  height='30' width='15.0'>\n")
			buffer.WriteString("<table border='5'  height='30' WIDTH='15.0' cellborder='1'>\n")
			buffer.WriteString(" <tr>  <td height='60' colspan='100%'>EXTENDIDA <br/>" + convertName(Particiones[i].Part_name[:]) + " <br/> Ocupado:" + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>  </tr>\n<tr>")
			var InicioExtendida int64 = Particiones[i].Part_start
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Part_next == -1 {
				fmt.Println("No hay particiones logicas")
			} else {
				var EspacioUtilizado int64 = 0
				cont := 0
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Part_next == -1 {
						break
					} else {
						f.Seek(ebr.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						EspacioUtilizado += ebr.Part_size
						PorcentajeUtilizado = (float64(ebr.Part_start) / float64(Particiones[i].Part_size)) * 100
						buffer.WriteString("<td height='30'>EBR</td><td height='30'> Logica:  " + convertName(ebr.Part_name[:]) + " " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
						cont++
					}
				}
				if (Particiones[i].Part_size - EspacioUtilizado) > 0 {
					PorcentajeUtilizado = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
					buffer.WriteString("<td height='30' width='100%'>Libre: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")
				}
			}
			buffer.WriteString("</tr>\n")
			buffer.WriteString("</table>\n</td>")
		}

	}
	if (TamanioDisco - EspacioUtilizado) > 0 {
		PorcentajeUtilizado = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
		buffer.WriteString("<td height='30' width='75.0'>Libre: " + strconv.Itoa(int(PorcentajeUtilizado)) + "%</td>")

	}
	buffer.WriteString("     </tr>\n</table>\n>];\n}")
	var datos string
	datos = string(buffer.String())
	fmt.Print(datos)
	CrearElReporte(path, datos)
	return false
}

func RecorrerListaDisco(id string, ListaDiscos *list.List) (string, string, string) {
	Id := strings.ReplaceAll(id, "19", "")
	//NoParticion := Id[1:]
	IdDisco := Id[1:2]
	pathDisco := ""
	nombreParticion := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco Disk
		disco = element.Value.(Disk)

		if BytesToString(disco.Id) == IdDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")
					nombreParticion = mountTemp.NombreParticion
					pathDisco = disco.Path
					nombreDisco = disco.NombreDisco
					return pathDisco, nombreParticion, nombreDisco

				}
			}

		}
		element.Value = disco
	}
	return "", "", ""
}

func convertName(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}
func convertBloqueData(c []byte) string {
	if c[0] == 32 {
		return " "
	}
	n := -1
	for i, b := range c {
		if b == 32 || b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func EstaLlenoDD(posicion int64, inicioDD int64, cantidadDD int64, pathDisco string) bool {
	estaLleno := false
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(inicioDD, 0)
	dd := DD{}
	for i := 0; i < int(cantidadDD); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			break
		}
		if i == int(posicion) {
			for j := 0; j < 5; j++ {
				if len(convertName(dd.Dd_array_files[j].Dd_file_nombre[:])) > 0 {
					estaLleno = true
					break
				} else {
					estaLleno = false
				}
			}
		}
	}
	return estaLleno
}
func CreateArchivo(path string, data string) {
	propiedades := strings.Split(path, "/")
	nombreArchivo := propiedades[len(propiedades)-1]
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
	executeComand("dot -Tpdf " + path + " -o " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	executeComand("xdg-open " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	executeComand("xdg-open " + path)

}

func CrearElReporte(path string, data string) {
	propiedades := strings.Split(path, "/")
	nombreArchivo := propiedades[len(propiedades)-1]
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
	executeComand("dot -Tpdf " + path + " -o " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	executeComand("xdg-open " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	executeComand("xdg-open " + path)

}

func GraficarTree(idParticion string, pathCarpeta string, ruta string, ListaDiscos *list.List) bool {
	var buffer bytes.Buffer
	buffer.WriteString("digraph grafica{\nrankdir=LR;\nnode [shape = record, style=filled, fillcolor=seashell2];\n")
	sb := SuperBloque{}
	var dos [15]byte
	avd := AVD{}
	var strArray [100]string
	//var InicioParticion int64 =0
	pathDisco, nombreParticion, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	sb, _ = DevolverSuperBloque(pathDisco, nombreParticion)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	/*
	   Graficar AVD's
	*/
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if avd.Avd_nombre_directotrio == dos {
			break
		}
		for j := 0; j < 6; j++ {
			if avd.Avd_ap_array_subdirectorios[j] != -1 {
				buffer.WriteString("nodo" + strconv.Itoa(i) + ":f" + strconv.Itoa(j) + " -> nodo" + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[j])) + "\n")
			} else {
				break
			}
		}
		if avd.Avd_ap_arbol_virtual_directorio != -1 {
			buffer.WriteString("nodo" + strconv.Itoa(i) + ":f7" + " -> nodo" + strconv.Itoa(int(avd.Avd_ap_arbol_virtual_directorio)) + "\n")
		}
		if EstaLlenoDD(avd.Avd_ap_detalle_directorio, sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
			strArray[i] = convertName(avd.Avd_nombre_directotrio[:])
			buffer.WriteString("nodo" + strconv.Itoa(i) + ":f6 -> node" + strconv.Itoa(int(avd.Avd_ap_detalle_directorio)) + "\n")
		}
		buffer.WriteString("nodo" + strconv.Itoa(i) + "[ shape=record, label =\"" + "{" + convertName(avd.Avd_nombre_directotrio[:]) + "|{<f0> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[0])) + "|<f1>" + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[1])) + "|<f2> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[2])) + "|<f3> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[3])) + "|<f4> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[4])) + "|<f5>" + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[5])) + "|<f6>" + strconv.Itoa(int(avd.Avd_ap_detalle_directorio)) + "|<f7> " + strconv.Itoa(int(avd.Avd_ap_arbol_virtual_directorio)) + "}}\"];\n")
	}
	/*
	   Graficar DD's
	*/
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	dd := DD{}
	for i := 0; i < int(sb.Sb_detalle_directorio_count); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			break
		}
		//fmt.Println(EstaLlenoDD(int64(i),sb.Sb_ap_detalle_directorio,sb.Sb_detalle_directorio_count,pathDisco),i)
		if EstaLlenoDD(int64(i), sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
			for j := 0; j < 5; j++ {
				if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) != convertName(dos[:]) {
					buffer.WriteString("node" + strconv.Itoa(i) + ":f" + strconv.Itoa(j+1) + "->  nodex" + strconv.Itoa(int(dd.Dd_array_files[j].Dd_file_ap_inodo)) + "\n")
				}
			}
			buffer.WriteString("node" + strconv.Itoa(i) + "[shape=record, label=\"" + "{ dd " + strArray[i] + "|")
			for j := 0; j < 5; j++ {
				if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) != convertName(dos[:]) {
					buffer.WriteString("{<f" + strconv.Itoa(j) + "> " + convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) + "| <f" + strconv.Itoa(j+1) + "> " + strconv.Itoa(int(dd.Dd_array_files[j].Dd_file_ap_inodo)) + "} |")
				} else {
					buffer.WriteString("{-1 | } |")
				}

			}
			if dd.Dd_ap_detalle_directorio != -1 {
				buffer.WriteString("{" + strconv.Itoa(int(dd.Dd_ap_detalle_directorio)) + " | <f10>  }}\"];\n")
				buffer.WriteString("node" + strconv.Itoa(i) + ":f10 -> " + "node" + strconv.Itoa(int(dd.Dd_ap_detalle_directorio)))
			} else {
				buffer.WriteString("{*1 | <f10>  }}\"];\n")
			}
			buffer.WriteString("\n")
		}
	}
	/*
	   Graficar Inodo's
	   X para identificarlos
	*/
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	inodo := INODO{}
	for i := 0; i < int(sb.S_inodes_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if inodo.I_count_inodo == -1 {
			break
		}
		if inodo.I_ao_indirecto != -1 {
			buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + "[shape=record, label=\"{Inodo" + strconv.Itoa(int(inodo.I_count_inodo)) + "|{" + strconv.Itoa(int(inodo.I_array_bloques[0])) + "| <f0> }|{" + strconv.Itoa(int(inodo.I_array_bloques[1])) + "| <f1> }|{" + strconv.Itoa(int(inodo.I_array_bloques[2])) + " | <f2> }|{" + strconv.Itoa(int(inodo.I_array_bloques[3])) + "| <f3> }|{" + strconv.Itoa(int(inodo.I_ao_indirecto)) + " | <f4> }}\"];" + "\n")
			buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f4 ->" + "nodex" + strconv.Itoa(int(inodo.I_ao_indirecto)) + "\n")
			for h := 0; h < 4; h++ {
				if inodo.I_array_bloques[h] == -1 {
					break
				} else {
					buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f" + strconv.Itoa(h) + "-> data" + strconv.Itoa(int(inodo.I_array_bloques[h])) + "\n")
				}
			}
		} else {
			buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + "[shape=record, label=\"{Inodo" + strconv.Itoa(int(inodo.I_count_inodo)) + "|{" + strconv.Itoa(int(inodo.I_array_bloques[0])) + "| <f0> }|{" + strconv.Itoa(int(inodo.I_array_bloques[1])) + "| <f1> }|{" + strconv.Itoa(int(inodo.I_array_bloques[2])) + " | <f2> }|{" + strconv.Itoa(int(inodo.I_array_bloques[3])) + "| <f3> }|{*" + strconv.Itoa(int(inodo.I_ao_indirecto)) + " | <f4> }}\"];" + "\n")
			for h := 0; h < 4; h++ {
				if inodo.I_array_bloques[h] == -1 {
					break
				} else {
					buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f" + strconv.Itoa(h) + "-> data" + strconv.Itoa(int(inodo.I_array_bloques[h])) + "\n")
				}
			}
		}
	}
	/*
	   Graficar Bloque's
	*/
	f.Seek(sb.Sb_ap_bloques, 0)
	data := BLOQUECAR{}
	for i := 0; i < int(sb.S_inodes_count); i++ {
		err = binary.Read(f, binary.BigEndian, &data)
		if data.B_content[0] == 0 {
			break
		}
		buffer.WriteString("data" + strconv.Itoa(i) + "[shape=record, label=\"{data| <f1> " + convertBloqueData(data.B_content[:]) + "}}\"];\n")

	}
	buffer.WriteString("\n}")
	var datos string
	datos = string(buffer.String())
	CrearElReporte(pathCarpeta, datos)
	return false
}

func Reportefile(idParticion string, pathCarpeta string, ruta string, ListaDiscos *list.List) bool {
	var bloquesGraficar [100]int
	carpetas_Graficar := strings.Split(ruta, "/")
	var buffer bytes.Buffer
	var noDirectorio int64 = 0
	buffer.WriteString("digraph grafica{\nrankdir=TB;\nnode [shape = record, style=filled, fillcolor=sienna1];\n")
	sb := SuperBloque{}
	var dos [15]byte
	avd := AVD{}
	var strArray [100]string
	//var InicioParticion int64 =0
	pathDisco, nombreParticion, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	sb, _ = DevolverSuperBloque(pathDisco, nombreParticion)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	f.Seek(sb.Sb_ap_arbol_directorio, 0)
	for i := 0; i < int(sb.Sb_arbol_virtual_count); i++ {
		err = binary.Read(f, binary.BigEndian, &avd)
		if avd.Avd_nombre_directotrio == dos {
			break
		}
		if convertName(avd.Avd_nombre_directotrio[:]) == carpetas_Graficar[len(carpetas_Graficar)-2] {
			for j := 0; j < 6; j++ {
				if avd.Avd_ap_array_subdirectorios[j] != -1 {
					buffer.WriteString("nodo" + strconv.Itoa(i) + ":f" + strconv.Itoa(j) + " -> nodo" + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[j])) + "\n")
				} else {
					break
				}
			}
			if avd.Avd_ap_arbol_virtual_directorio != -1 {
				buffer.WriteString("nodo" + strconv.Itoa(i) + ":f7" + " -> nodo" + strconv.Itoa(int(avd.Avd_ap_arbol_virtual_directorio)) + "\n")
			}
			noDirectorio = avd.Avd_ap_detalle_directorio
			if EstaLlenoDD(avd.Avd_ap_detalle_directorio, sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
				strArray[i] = convertName(avd.Avd_nombre_directotrio[:])
				buffer.WriteString("nodo" + strconv.Itoa(i) + ":f6 -> node" + strconv.Itoa(int(avd.Avd_ap_detalle_directorio)) + "\n")
			}
			buffer.WriteString("nodo" + strconv.Itoa(i) + "[ shape=record, label =\"" + "{" + convertName(avd.Avd_nombre_directotrio[:]) + "|{<f0> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[0])) + "|<f1>" + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[1])) + "|<f2> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[2])) + "|<f3> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[3])) + "|<f4> " + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[4])) + "|<f5>" + strconv.Itoa(int(avd.Avd_ap_array_subdirectorios[5])) + "|<f6>" + strconv.Itoa(int(avd.Avd_ap_detalle_directorio)) + "|<f7> " + strconv.Itoa(int(avd.Avd_ap_arbol_virtual_directorio)) + "}}\"];\n")
		}
	}
	/*
	   Graficar DD's
	*/
	noInodoGraficar := 0
	f.Seek(sb.Sb_ap_detalle_directorio, 0)
	dd := DD{}
	for i := 0; i < int(sb.Sb_detalle_directorio_count); i++ {
		err = binary.Read(f, binary.BigEndian, &dd)
		if dd.Ocupado == 0 {
			break
		}
		if noDirectorio == int64(i) {
			if EstaLlenoDD(int64(i), sb.Sb_ap_detalle_directorio, sb.Sb_detalle_directorio_count, pathDisco) {
				for j := 0; j < 5; j++ {
					if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) == carpetas_Graficar[len(carpetas_Graficar)-1] {
						noInodoGraficar = int(dd.Dd_array_files[j].Dd_file_ap_inodo)
						buffer.WriteString("node" + strconv.Itoa(i) + ":f" + strconv.Itoa(j+1) + "->  nodex" + strconv.Itoa(int(dd.Dd_array_files[j].Dd_file_ap_inodo)) + "\n")
					}
				}
				buffer.WriteString("node" + strconv.Itoa(i) + "[shape=record, label=\"" + "{ dd " + strArray[i] + "|")
				for j := 0; j < 5; j++ {
					if convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) != convertName(dos[:]) {
						buffer.WriteString("{<f" + strconv.Itoa(j) + "> " + convertName(dd.Dd_array_files[j].Dd_file_nombre[:]) + "| <f" + strconv.Itoa(j+1) + "> " + strconv.Itoa(int(dd.Dd_array_files[j].Dd_file_ap_inodo)) + "} |")
					} else {
						buffer.WriteString("{-1 | } |")
					}

				}
				if dd.Dd_ap_detalle_directorio != -1 {
					noDirectorio = dd.Dd_ap_detalle_directorio
					buffer.WriteString("{" + strconv.Itoa(int(dd.Dd_ap_detalle_directorio)) + " | <f10>  }}\"];\n")
					buffer.WriteString("node" + strconv.Itoa(i) + ":f10 -> " + "node" + strconv.Itoa(int(dd.Dd_ap_detalle_directorio)))
				} else {
					buffer.WriteString("{*1 | <f10>  }}\"];\n")
				}
				buffer.WriteString("\n")
			}
		}
	}
	/*
	   Graficar Inodo's
	   X para identificarlos
	*/
	f.Seek(sb.Sb_ap_tabla_inodo, 0)
	inodo := INODO{}
	cont1 := 0
	for i := 0; i < int(sb.S_inodes_count); i++ {
		err = binary.Read(f, binary.BigEndian, &inodo)
		if inodo.I_count_inodo == -1 {
			break
		}
		if noInodoGraficar == i {
			if inodo.I_ao_indirecto != -1 {
				noInodoGraficar = i + 1
				buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + "[shape=record, label=\"{Inodo" + strconv.Itoa(int(inodo.I_count_inodo)) + "|{" + strconv.Itoa(int(inodo.I_array_bloques[0])) + "| <f0> }|{" + strconv.Itoa(int(inodo.I_array_bloques[1])) + "| <f1> }|{" + strconv.Itoa(int(inodo.I_array_bloques[2])) + " | <f2> }|{" + strconv.Itoa(int(inodo.I_array_bloques[3])) + "| <f3> }|{" + strconv.Itoa(int(inodo.I_ao_indirecto)) + " | <f4> }}\"];" + "\n")
				buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f4 ->" + "nodex" + strconv.Itoa(int(inodo.I_ao_indirecto)) + "\n")
				for h := 0; h < 4; h++ {
					if inodo.I_array_bloques[h] == -1 {
						break
					} else {
						buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f" + strconv.Itoa(h) + "-> data" + strconv.Itoa(int(inodo.I_array_bloques[h])) + "\n")
						bloquesGraficar[cont1] = int(inodo.I_array_bloques[h])
						cont1++
					}
				}
			} else {
				buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + "[shape=record, label=\"{Inodo" + strconv.Itoa(int(inodo.I_count_inodo)) + "|{" + strconv.Itoa(int(inodo.I_array_bloques[0])) + "| <f0> }|{" + strconv.Itoa(int(inodo.I_array_bloques[1])) + "| <f1> }|{" + strconv.Itoa(int(inodo.I_array_bloques[2])) + " | <f2> }|{" + strconv.Itoa(int(inodo.I_array_bloques[3])) + "| <f3> }|{*" + strconv.Itoa(int(inodo.I_ao_indirecto)) + " | <f4> }}\"];" + "\n")
				for h := 0; h < 4; h++ {
					if inodo.I_array_bloques[h] == -1 {
						break
					} else {
						bloquesGraficar[cont1] = int(inodo.I_array_bloques[h])
						buffer.WriteString("nodex" + strconv.Itoa(int(inodo.I_count_inodo)) + " :f" + strconv.Itoa(h) + "-> data" + strconv.Itoa(int(inodo.I_array_bloques[h])) + "\n")
						cont1++
					}
				}
			}
		}
	}
	cont1 = 0
	f.Seek(sb.Sb_ap_bloques, 0)
	data := BLOQUECAR{}
	for i := 0; i < int(sb.S_bloques_count); i++ {
		err = binary.Read(f, binary.BigEndian, &data)
		if data.B_content[0] == 0 {
			break
		}
		if bloquesGraficar[cont1] == i {
			buffer.WriteString("data" + strconv.Itoa(i) + "[shape=record, label=\"{data| <f1> " + convertBloqueData(data.B_content[:]) + "}}\"];\n")
			cont1++
		}
	}
	//Crear Archivo
	buffer.WriteString("\n}")
	var datos string
	datos = string(buffer.String())
	CreateArchivo(pathCarpeta, datos)
	return false
}
