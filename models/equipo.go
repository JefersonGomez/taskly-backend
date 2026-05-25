package models

import "gorm.io/gorm"

//Equipo   { ID, Nombre, Descripcion, OwnerID }

type Equipo struct {
	gorm.Model
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	OwnerID     uint      `json:"ownerid"`
	Tareas      []Tarea   `json:"tareas" gomr:"foreignKey:EquipoID"`
	Miembros    []Miembro `json:"miembros" gorm:"foreignKey:EquipoID"`
}
