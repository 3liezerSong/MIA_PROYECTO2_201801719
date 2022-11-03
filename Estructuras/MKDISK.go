package Estructuras

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func EjecutarComandoMKDISK(nombreComando string, propiedadesTemp []Propiedad, cont int) (ParamValidos bool) {
	dt := time.Now()
	mbr1 := MBR{}
	copy(mbr1.Mbr_fecha_creacion[:], dt.String())
	mbr1.Mbr_dsk_signature = int64(rand.Intn(100) + cont)
	var buffer [1024]byte

	fmt.Println("############# Comando MKDISK #############")
	ParamValidos = true
	pathwithoutfile := ""
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {

		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-size":
				propiedades[0] = propiedadTemp.Val
			case "-unit":
				propiedades[2] = strings.ToLower(propiedadTemp.Val)
			case "-path":
				propiedades[3] = propiedadTemp.Val
				res1 := strings.Split(propiedadTemp.Val, "/")
				pathwithoutfile = propiedadTemp.Val[0 : len(propiedadTemp.Val)-len(res1[len(res1)-1])-1]
				crearDirectorioSiNoExiste(pathwithoutfile)

				f, err := os.Create(pathwithoutfile + "/" + res1[len(res1)-1])
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}

		tamaniototal, _ := strconv.ParseInt(propiedades[0], 10, 64)
		for i := 0; i < 1024; i++ {
			buffer[i] = '0'
		}
		if propiedades[2] == "k" {
			mbr1.Mbr_tamano = ((tamaniototal * 1024) / 1024)

		} else {
			mbr1.Mbr_tamano = (tamaniototal * 1024 * 1024) / 1024
		}

		//Aquí inicializo las particiones
		for i := 0; i < 4; i++ {
			copy(mbr1.Particiones[i].Part_status[:], "0")
			copy(mbr1.Particiones[i].Part_type[:], "0")
			copy(mbr1.Particiones[i].Part_fit[:], "00")
			mbr1.Particiones[i].Part_start = 0
			mbr1.Particiones[i].Part_size = 0
			copy(mbr1.Particiones[i].Part_name[:], "0000000000000000")
		}
		//Escribo el MBR
		f, err := os.OpenFile(propiedades[3], os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		f.Seek(0, 0)

		err = binary.Write(f, binary.BigEndian, mbr1)
		for i := 0; i < int(mbr1.Mbr_tamano); i++ {
			err := binary.Write(f, binary.BigEndian, buffer)
			if err != nil {
				log.Fatalln(err, mbr1.Mbr_tamano)
			}
		}
		if err != nil {
			log.Fatalln(err, propiedades[3])
		}
		fmt.Println("Disco Creado Exitosamente")

		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func crearDirectorioSiNoExiste(directorio string) {
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.MkdirAll(directorio, 0755)
		if err != nil {
			println("No pude crear el directorio :c")
			panic(err)
		}
	}
}

func BytesToString(data [1]byte) string {
	return string(data[:])
}

func executeComand(comandos string) {
	args := strings.Split(comandos, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.CombinedOutput()
}
