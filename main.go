package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"

	L "MIA_PROYECTO2_201801719/Estructuras"

	Lectura "MIA_PROYECTO2_201801719/Estructuras"
)

func main() {
	fmt.Println("\033[2J")
	ListaDiscos := list.New()
	LlenarListaDisco(ListaDiscos)
	var comando string = ""
	escaner := bufio.NewScanner(os.Stdin)
	for comando != "exit" {
		fmt.Println("--------------------------------------------------")
		fmt.Println("|---------Eliezer Abraham Zapeta Alvarado--------|")
		fmt.Println("|----------------- 201801719 --------------------|")
		fmt.Println("--------------------------------------------------")
		fmt.Println("POR FAVOR INGRESE EL COMANDO A EJECUTAR")
		escaner.Scan()
		comando = escaner.Text()
		if comando != "" {
			Lectura.AnalizarComando(comando, ListaDiscos)
		}

	}

}

func LlenarListaDisco(ListaDiscos *list.List) {
	IdDisco := [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I",
		"J", "K", "L", "M", "N", "O", "P", "Q",
		"R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for i := 0; i < 26; i++ {
		disco := L.Disk{}
		copy(disco.Estado[:], "0")
		copy(disco.Id[:], IdDisco[i])
		for j := 0; j < len(disco.Particiones); j++ {
			mount := L.MOUNT{}
			mount.NombreParticion = ""
			mount.Id = strconv.Itoa(j + 1)
			copy(mount.Estado[:], "0")
			disco.Particiones[j] = mount
		}
		ListaDiscos.PushBack(disco)
	}
}
