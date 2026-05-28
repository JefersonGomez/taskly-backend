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

// TareaConAsignadoResponse Tarea con nombre del usuario asignado
// @Description Respuesta de tarea que incluye el nombre del miembro asignado (si existe)
type TareaConAsignadoResponse struct {
	ID             uint   `json:"id" example:"1"`
	Titulo         string `json:"titulo" example:"Fix login bug"`
	Descripcion    string `json:"descripcion" example:"Corregir error en autenticación"`
	Estado         string `json:"estado" example:"en_progreso"`
	EquipoID       uint   `json:"equipo_id" example:"3"`
	AsignadoID     uint   `json:"asignado_id,omitempty" example:"5"`
	NombreAsignado string `json:"nombreAsignado,omitempty" example:"Ana López"`
	CreatedAt      string `json:"created_at" example:"2026-05-28T10:00:00Z"`
}
