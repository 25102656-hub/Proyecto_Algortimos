package modelos

type Lugar struct {
	Id        int     `json:"id"`
	Nombre    string  `json:"nombre"`
	Categoria string  `json:"categoria"`
	Ciudad    string  `json:"ciudad"`
	Latitud   float64 `json:"latitud"`
	Longitud  float64 `json:"longitud"`
}