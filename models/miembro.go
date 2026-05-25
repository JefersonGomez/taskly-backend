package models

import "gorm.io/gorm"

//Miembro  { ID, EquipoID, UsuarioID, Rol }

type Miembro struct {
	gorm.Model
	EquipoID  uint    `json:"equipoid"`
	UsuarioID uint    `json:"usuarioid"`
	Rol       string  `json:"rol"`
	Usuario   Usuario `json:"usuario" gorm:"foreignKey:usuarioID"`
	Equipo    Equipo  `json:"equipo" gorm:"foreignKey:equipoID"`
}
