package Estructuras

import (
	"container/list"
	"fmt"
	"strings"
)

//EStruc de las particiones montadas

var global string = ""
var global2 string = ""

//Aqui voy a leer y reconocer los comandos

func AnalizarComando(data string, ListaDiscos *list.List) {
	ListaComandos := list.New()
	lineaComando := strings.Split(data, "\n") //hago split el salto de linea

	var c Comando
	for i := 0; i < len(lineaComando); i++ {
		Comment := lineaComando[i][0:1]
		if Comment != "#" {
			comando := lineaComando[i]

			if strings.Contains(lineaComando[i], "\\*") {
				comando = strings.Replace(lineaComando[i], "\\*", " ", 1) + lineaComando[i+1]
				i = i + 1
			}
			propiedades := strings.Split(string(comando), " ")
			//Aqui va el identificador del comando
			nombreComando := propiedades[0]
			//Estructura para el comando
			c = Comando{Name: strings.ToLower(nombreComando)}
			Tempropiedades := make([]Propiedad, len(propiedades)-1)
			for i := 1; i < len(propiedades); i++ {
				if propiedades[i] == "" {
					continue
				} else if propiedades[i] == "-p" {
					Tempropiedades[i-1] = Propiedad{Name: "-p",
						Val: "-p"}
				} else {
					if strings.Contains(propiedades[i], "=") {
						valor_propiedad_Comando := strings.Split(propiedades[i], "=")

						Tempropiedades[i-1] = Propiedad{Name: valor_propiedad_Comando[0],
							Val: valor_propiedad_Comando[1]}
					} else {
						Tempropiedades[i-1] = Propiedad{Name: "-sigue",
							Val: propiedades[i]}
					}
				}
			}
			c.Propiedades = Tempropiedades
			//Agregando el comando a la lista comandos
			ListaComandos.PushBack(c)
		} else {
			comando := lineaComando[i]
			fmt.Println(comando)
		}
	}
	RecorrerListaComando(ListaComandos, ListaDiscos)
}

func RecorrerListaComando(ListaComandos *list.List, ListaDiscos *list.List) {
	var ParamValidos bool = true
	var cont = 1
	for element := ListaComandos.Front(); element != nil; element = element.Next() {
		var comandoTemp Comando
		comandoTemp = element.Value.(Comando)
		//Lista de propiedades del Comando
		switch strings.ToLower(comandoTemp.Name) {
		case "mkdisk":
			ParamValidos = EjecutarComandoMKDISK(comandoTemp.Name, comandoTemp.Propiedades, cont)
			cont++

			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "rmdisk":
			ParamValidos = EjecutarComandoRMDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "fdisk":
			ParamValidos = EjecutarComandoFDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if !ParamValidos {
			}
		case "mount":
			if len(comandoTemp.Propiedades) != 0 {
				ParamValidos = EjejcutarComandoMount(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
				if !ParamValidos {
					fmt.Println("Parametros Invalidos")
				}
			} else {
				fmt.Println("hola") //	EjecutarReporteMount(ListaDiscos)
			}

		case "exit":
			fmt.Println("Finalizo la Ejecucion")
		case "pause":
			fmt.Println("Presione una tecla para Continuar")
			fmt.Scanln()
		case "exec":
			ParamValidos = EjecutarComandoExec(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)

			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "mkdir":
			ParamValidos = EjecutarComandoMKDIR(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "mkfile":
			ParamValidos = EjecutarComandoMKFILE(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "mkfs":
			ParamValidos = EjecutarComandoMKFS(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "rep":
			ParamValidos = EjecutarComandoRep(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "rmgrp":
			ParamValidos = EjecutarComandoMKGRP(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "login":
			ParamValidos, global, global2 = EjecutarComandoLogin(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if !ParamValidos {
				fmt.Println("Parametros Invalidos")
			}
		case "logout":
			if global == "" {
				fmt.Println("Debe Iniciar Sesion")
			} else {
				global = ""
				global2 = ""
				fmt.Println("Se ha cerrado la sesión")

			}
		default:
			fmt.Println("Error al Ejecutar el Comando")
		}
	}
}
