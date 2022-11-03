package Estructuras

type Propiedad struct {
	Name string
	Val  string
}

type Comando struct {
	Name        string
	Propiedades []Propiedad
}

//Estruct para el MBR
type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int64
	// dsk_fit            [1]byte
	Particiones [4]Particion
}

type MOUNT struct {
	NombreParticion string
	Id              string
	Estado          [1]byte
	EstadoMKS       [1]byte
}

//Estruct las Particiones
type Particion struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [2]byte
	Part_start  int64
	Part_size   int64
	Part_name   [15]byte
}

//Estruct Disco
type Disk struct {
	NombreDisco string
	Path        string
	Id          [1]byte
	Estado      [1]byte
	Particiones [100]MOUNT
}

//Estruct para el EBR
type EBR struct {
	Part_status [1]byte
	Part_fit    [2]byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [15]byte
}

//Para apuntar
type Bitacora struct {
	Log_tipo_operacion [19]byte
	Log_tipo           [1]byte
	Log_nombre         [35]byte
	Log_Contenido      [25]byte
	Log_fecha          [19]byte
	Size               int64
}

//SuperBloque
type SuperBloque struct {
	Sb_nombre_hd                         [15]byte
	Sb_arbol_virtual_count               int64
	Sb_detalle_directorio_count          int64
	S_inodes_count                       int64
	S_bloques_count                      int64
	Sb_arbol_virtual_free                int64
	Sb_detalle_directorio_free           int64
	Sb_inodos_free                       int64
	Sb_bloques_free                      int64
	Sb_date_creacion                     [19]byte
	Sb_date_ultimo_montaje               [19]byte
	Sb_montajes_count                    int64
	Sb_ap_bitmap_arbol_directorio        int64
	Sb_ap_arbol_directorio               int64
	Sb_ap_bitmap_detalle_directorio      int64
	Sb_ap_detalle_directorio             int64
	Sb_ap_bitmap_tabla_inodo             int64
	Sb_ap_tabla_inodo                    int64
	Sb_ap_bitmap_bloques                 int64
	Sb_ap_bloques                        int64
	Sb_ap_log                            int64
	Sb_size_struct_arbol_directorio      int64
	Sb_size_struct_Detalle_directorio    int64
	Sb_size_struct_inodo                 int64
	Sb_size_struct_bloque                int64
	Sb_first_free_bit_arbol_directorio   int64
	Sb_first_free_bit_detalle_directorio int64
	Sb_dirst_free_bit_tabla_inodo        int64
	Sb_first_free_bit_bloques            int64
	Sb_magic_num                         int64
	InicioCopiaSB                        int64
	ConteoAVD                            int64
	ConteoDD                             int64
	ConteoInodo                          int64
	ConteoBloque                         int64
}

type ArregloDD struct {
	Dd_file_nombre            [15]byte
	Dd_file_ap_inodo          int64
	Dd_file_date_creacion     [19]byte
	Dd_file_date_modificacion [19]byte
}
type DD struct {
	Dd_array_files           [5]ArregloDD
	Dd_ap_detalle_directorio int64
	Ocupado                  int8
}
type AVD struct {
	Avd_fecha_creacion              [19]byte
	Avd_nombre_directotrio          [15]byte
	Avd_ap_array_subdirectorios     [6]int64
	Avd_ap_detalle_directorio       int64
	Avd_ap_arbol_virtual_directorio int64
	Avd_proper                      [10]byte
}

type INODO struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64
	I_array_bloques           [4]int64
	I_ao_indirecto            int64
	I_id_proper               int64
}

type BLOQUECAR struct {
	B_content [25]byte
}

func BytesNombreParticion(data [15]byte) string {
	return string(data[:])
}

func ConvertirData(data [25]byte) string {
	return string(data[:])
}
