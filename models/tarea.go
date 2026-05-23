package models

import "gorm.io/gorm"

//Tarea    { ID, Titulo, Descripcion, Estado,
//           Prioridad, FechaLimite, EquipoID, AsignadoID }

type Tarea struct {
	gorm.Model
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Estado      string `json:"estado"`
	Prioridad   string `json:"prioridad"`
	FechaLimite string `json:"fechaLimite"`
	EquipoID    uint   `json:"equipoid"`
	AsignadoID  uint   `json:"asignadoid"`
}
